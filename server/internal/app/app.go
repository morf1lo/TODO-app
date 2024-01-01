package app

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/morf1lo/TODO-app/internal/db"
	"github.com/morf1lo/TODO-app/internal/routes"
	"github.com/morf1lo/TODO-app/config"
)

func App() {
	config.Init()

	client, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.New()
	router.Use(cors.Default())

	router.SetTrustedProxies(nil)

	routes.SetupUserRoutes(router, client.Database("TODO").Collection("users"))

	router.Run(":8080")
}
