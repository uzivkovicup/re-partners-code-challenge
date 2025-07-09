package constants

import "time"

type AppEnv string

const (
	ProductionEnv AppEnv = "production"
	StageEnv      AppEnv = "stage"
	DevelopEnv    AppEnv = "develop"
	LocalEnv      AppEnv = "local"
)

const ReadHeaderTimeout = 10 * time.Second

const ShutdownTimeout = 10 * time.Second

const GormLoggerSlowThreshold = 200 * time.Millisecond
