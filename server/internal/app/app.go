package app

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func App() {
	router := gin.New()
	router.Use(cors.Default())

	router.SetTrustedProxies(nil)

	router.GET("/api/greeting/:name", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": fmt.Sprintf("Hello, %s!", c.Param("name"))})
	})

	router.Run(":8080")
}
