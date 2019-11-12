package main

import (
	"io"
	"net/http"
	"os"
	"rest-api-gin-mongo/controllers"
	"rest-api-gin-mongo/db"

	"github.com/gin-gonic/gin"
)

var (
	config = db.Config{}
)

func init() {
	gin.DisableConsoleColor()
	config.Read("config.toml")

	// Logging to a file.
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
}

func SetupRouter() *gin.Engine {
	gin.SetMode(config.Gin_mode)
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

	router.GET("ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return router
}

func main() {
	router := SetupRouter()
	router.Run(":" + config.Gin_port)
}
