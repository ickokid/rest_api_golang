package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"rest-api-gin-mongo/controllers"
	"rest-api-gin-mongo/db"
	"rest-api-gin-mongo/utils"

	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-gonic/gin"
)

var (
	config = db.Config{}
)

func init() {
	gin.DisableConsoleColor()
	config.Read("config.toml")

	// Logging to a file.
	ginLog, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(ginLog)
}

func TokenAuthMiddleware() gin.HandlerFunc {
	requiredKey := config.Key

	if requiredKey == "" {
		log.Fatal("Please set API Key environment variable")
	}

	return func(c *gin.Context) {
		key := c.Request.FormValue("key")

		if key == "" {
			utils.Forbidden(c)
			return
		}

		if key != requiredKey {
			utils.Forbidden(c)
			return
		}

		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	gin.SetMode(config.Gin_mode)
	router := gin.Default()
	router.Use(limit.MaxAllowed(20))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	/*v1 := router.Group("api/v1")
	{
		v1.GET("/instructions", controllers.GetInstructions)
		v1.GET("/instructions/:id", controllers.GetInstruction)
		v1.POST("/instructions", controllers.PostInstruction)
		v1.PUT("/instructions/:id", controllers.UpdateInstruction)
		v1.DELETE("/instructions/:id", controllers.DeleteInstruction)
	}*/

	v1 := router.Group("api/v1")
	v1.Use(TokenAuthMiddleware())
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
