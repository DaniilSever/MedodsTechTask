package main

import (
	"github.com/MedodsTechTask/app/core"
	"github.com/MedodsTechTask/app/user/auth"
	_ "github.com/MedodsTechTask/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Service API
// @version 1.0
// @description REST API for authentication
// @host localhost:8080
// @BasePath /api/v1
func main() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Инициализация роутов
	api := r.Group(core.BasePath)
	{
		auth.SetupRoutes(api.Group(core.UserAuthPath))
	}

	r.Run(":8080")
}
