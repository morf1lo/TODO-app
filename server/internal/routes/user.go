package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/morf1lo/TODO-app/internal/handlers"
)

func SetupUserRoutes(router *gin.Engine, collection *mongo.Collection) {
	router.POST("/api/users/create", handlers.CreateUser(collection))
	router.POST("/api/users/login", handlers.Login(collection))
}
