package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"rest-api-gin-mongo/controllers"
)

func SetupRouter() *gin.Engine {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	GIN_MODE := os.Getenv("GIN_MODE")

	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(GIN_MODE)
	gin.DisableConsoleColor()
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("api/v1")
	{
		v1.GET("/instructions", controllers.GetInstructions)
		v1.GET("/instructions/:id", controllers.GetInstruction)
		v1.POST("/instructions", controllers.PostInstruction)
		v1.PUT("/instructions/:id", controllers.UpdateInstruction)
		v1.DELETE("/instructions/:id", controllers.DeleteInstruction)
	}

	return router
}

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	router := SetupRouter()
	router.Run(":" + PORT)
}
