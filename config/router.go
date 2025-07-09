package config

import "github.com/gin-gonic/gin"

var DefaultRouter *gin.Engine

func InitRouter() *gin.Engine {
	if DefaultRouter == nil {
		DefaultRouter = gin.Default()
	}
	return DefaultRouter
}
