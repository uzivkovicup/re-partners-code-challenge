package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go-pack-calculator/config"
	"go-pack-calculator/db"
	_ "go-pack-calculator/docs" // Import swagger docs
	"go-pack-calculator/internal/adapters/primary/rest"
	"go-pack-calculator/internal/adapters/secondary/inmemory"
	"go-pack-calculator/internal/adapters/secondary/postgres"
	"go-pack-calculator/internal/application/services"
	"go-pack-calculator/internal/ports/secondary"
)

// @title           Pack Calculator API
// @version         1.0
// @description     API for calculating optimal packs for orders
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.basic  BasicAuth
func main() {
	// Get current dir
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	// Load config with Viper (will check both .env file and environment variables)
	cfg, err := config.LoadConfig(currentDir)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	r := config.InitRouter()

	// Create global variables
	var (
		protocol           = cfg.Protocol
		addr               = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		postgresDBHost     = cfg.PostgresDBHost
		postgresDBPort     = cfg.PostgresDBPort
		postgresDBUser     = cfg.PostgresDBUser
		postgresDBPassword = cfg.PostgresDBPassword
		postgresDBName     = cfg.PostgresDBName
		postgresDBSSLMode  = cfg.PostgresDBSSLMode
	)

	// Initialize repository
	var packSizeRepository secondary.PackSizeRepository

	// Connect to PostgresDB in production, use in-memory repository in test
	if cfg.Environment == "test" {
		log.Println("Using in-memory repository for testing")
		packSizeRepository = inmemory.NewPackSizeRepository()
	} else {
		// Connect to PostgresDB
		err = db.NewPostgresDB(
			postgresDBHost,
			postgresDBPort,
			postgresDBUser,
			postgresDBPassword,
			postgresDBName,
			postgresDBSSLMode,
		).Connect()
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		// Create repository using GORM
		packSizeRepository = postgres.NewPackSizeRepository(db.PostgresDB)
	}

	// Initialize application service
	packCalculatorService := services.NewPackCalculatorService(packSizeRepository)

	// Initialize REST handler
	packCalculatorHandler := rest.NewPackCalculatorHandler(packCalculatorService, packCalculatorService)

	// Register REST API routes
	packCalculatorHandler.RegisterRoutes(r)

	// Serve static files
	r.Static("/static", "./static")

	// Serve pack calculator UI
	r.GET("/", func(c *gin.Context) {
		c.File("static/packcalculator.html")
	})

	// Serve Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint for deployment monitoring
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Create a new server with the router
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Create a channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Server started on %s://%s\n", protocol, addr)
		log.Printf("REST API available at %s://%s/api\n", protocol, addr)
		log.Printf("Swagger documentation available at %s://%s/swagger/index.html\n", protocol, addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Block until we receive a signal
	<-stop
	log.Println("Shutting down server...")

	// Create a deadline context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
