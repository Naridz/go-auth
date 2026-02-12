package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	Auth "github.com/go-auth/auth"
	"github.com/go-auth/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDB()

	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/register", Auth.Register)
	r.POST("/login", Auth.Login)
	r.Run(":8080")
}
