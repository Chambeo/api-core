package main

import (
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
	userService := service.NewUser(userRepository)
	// Handler
	userHandler := handler.NewUserHandler(userService)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// TODO segurizar endpoints q apliquen
	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("/", userHandler.Create)
			users.GET("/:id", userHandler.Get)
			users.GET("/email/:email", userHandler.GetByEmail)
			users.PUT("/", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}

	}

	err = r.Run(":8080")
	if err != nil {
		panic(">>>>>>>>>>>>>> Error trying to start the application <<<<<<<<<<<<<<")
	}
}
