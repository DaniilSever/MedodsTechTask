package main

import (
	"fmt"

	"github.com/MedodsTechTask/app/core"
	"github.com/MedodsTechTask/app/user/auth"
	"github.com/MedodsTechTask/app/user/auth/configs"
	"github.com/MedodsTechTask/app/user/auth/repo"
	_ "github.com/MedodsTechTask/docs"
	"github.com/gin-contrib/cors"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Или конкретные домены
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	authRepo, err := repo.NewAuthRepo(configs.GetConfig())
	if err != nil {
		fmt.Printf("%s", err)
	}
	authAPI := auth.NewAPI(auth.NewAuthUseCase(authRepo))

	// Инициализация роутов
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := r.Group(core.BasePath)
	{
		authAPI.SetupRoutes(api.Group(core.UserAuthPath))
	}

	r.Run(":8080")
}
