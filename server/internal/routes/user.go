package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/morf1lo/TODO-app/internal/controllers"
)

func SetupUserRoutes(router *gin.Engine, collection *mongo.Collection) {
	router.POST("/api/users/create", controllers.CreateUser(collection))
	router.POST("/api/users/login", controllers.Login(collection))
}
