package main

import (
	authHandler "chambeo-api-core/internal/auth/handler"
	authService "chambeo-api-core/internal/auth/service"
	userHandler "chambeo-api-core/internal/users/handler"
	userRepository "chambeo-api-core/internal/users/repository"
	userService "chambeo-api-core/internal/users/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// DB

	//db, err := gorm.Open(postgres.Open("jdbc:postgresql://127.0.0.1:5432/chambeo"), &gorm.Config{}) // TODO

	dsn := "host=127.0.0.1 user=chambeo password=chambeo dbname=chambeo port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Repo
	usrRepository := userRepository.NewUser(*db)
	// Service
	authenticationService := authService.NewJWTService()
	usrService := userService.NewUser(usrRepository)
	// Handler
	usrHandler := userHandler.NewUserHandler(usrService)
	authenticationHandler := authHandler.NewAuthHandler(&authenticationService, usrService)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// TODO segurizar endpoints q apliquen
	v1 := r.Group("/api/v1")
	{
		usersRouting := v1.Group("/users")
		{
			usersRouting.POST("/", usrHandler.Create)
			usersRouting.GET("/:id", usrHandler.Get)
			usersRouting.GET("/email/:email", usrHandler.GetByEmail)
			usersRouting.PUT("/", usrHandler.Update)
			usersRouting.DELETE("/:id", usrHandler.Delete)
		}

		authRouting := v1.Group("/auth")
		{
			authRouting.POST("/token", authenticationHandler.GenerateToken)
			authRouting.GET("/token/validate", authenticationHandler.ValidateToken)
			authRouting.POST("/token/refresh", authenticationHandler.RefreshToken)
		}

	}

	err = r.Run(":8080")
	if err != nil {
		panic(">>>>>>>>>>>>>> Error trying to start the application <<<<<<<<<<<<<<")
	}
}
