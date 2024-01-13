package main

import (
	authHandler "chambeo-api-core/internal/auth/handler"
	auth "chambeo-api-core/internal/auth/service"
	"chambeo-api-core/internal/users/handler"
	"chambeo-api-core/internal/users/repository"
	"chambeo-api-core/internal/users/service"
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
	userRepository := repository.NewUser(*db)
	// Service
	authService := auth.NewJWTService()
	userService := service.NewUser(userRepository)
	// Handler
	userHandling := handler.NewUserHandler(userService)
	authHandling := authHandler.NewAuthHandler(&authService, userService)

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
			usersRouting.POST("/", userHandling.Create)
			usersRouting.GET("/:id", userHandling.Get)
			usersRouting.GET("/email/:email", userHandling.GetByEmail)
			usersRouting.PUT("/", userHandling.Update)
			usersRouting.DELETE("/:id", userHandling.Delete)
		}

		authRouting := v1.Group("/auth")
		{
			authRouting.POST("/token", authHandling.GenerateToken)
			authRouting.GET("/token/validate", authHandling.ValidateToken)
			authRouting.POST("/token/refresh", authHandling.RefreshToken)
		}

	}

	err = r.Run(":8080")
	if err != nil {
		panic(">>>>>>>>>>>>>> Error trying to start the application <<<<<<<<<<<<<<")
	}
}
