package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/morf1lo/TODO-app/internal/handlers"
	"github.com/morf1lo/TODO-app/internal/middlewares"
)

func SetupTodoRoutes(router *gin.Engine, collection *mongo.Collection) {
	router.POST("/api/todos/create", middlewares.AuthMiddleware(), handlers.CreateTodo(collection))
	router.GET("/api/todos", middlewares.AuthMiddleware(), handlers.FindAllTodos(collection))
	router.PUT("/api/todos/update/:id", middlewares.AuthMiddleware(), handlers.UpdateTodo(collection))
	router.DELETE("/api/todos/delete/:id", middlewares.AuthMiddleware(), handlers.DeleteTodo(collection))
}
