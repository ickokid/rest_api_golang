package main

import (
	"log"
	"net/http"
	"os"
	"rest-api-gin-mysql/controllers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
		instruct := new(controllers.InstructionController)
		v1.GET("/instructions", instruct.GetInstructions)
		v1.GET("/instructions/:id", instruct.GetInstruction)
		v1.POST("/instructions", instruct.PostInstruction)
		v1.PUT("/instructions/:id", instruct.UpdateInstruction)
		v1.DELETE("/instructions/:id", instruct.DeleteInstruction)
	}

	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})

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
