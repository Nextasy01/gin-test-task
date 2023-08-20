package main

import (
	"os"

	"github.com/Nextasy01/gin-test-task/middleware"
	"github.com/Nextasy01/gin-test-task/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {

	server := gin.Default()

	public := server.Group("/")
	routes.PublicRoutes(public)

	private := server.Group("/")
	private.Use(middleware.AuthorizeJWT())
	routes.PrivateRoutes(private)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	server.Run(":" + port)

}
