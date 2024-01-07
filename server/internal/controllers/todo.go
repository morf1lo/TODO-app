package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/morf1lo/TODO-app/internal/models"
)

func CreateTodo(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString("uname")
		if len(username) <= 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not authorized"})
		}

		var todo models.Todo
		
		todo.Username = username

		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := todo.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := collection.InsertOne(context.TODO(), todo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func FindAllTodos(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString("uname")
		if len(username) <= 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not authorized"})
			return
		}

		cursor, err := collection.Find(context.TODO(), bson.M{"username": username})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		defer cursor.Close(context.TODO())

		var todos []models.Todo

		for cursor.Next(context.TODO()) {
			var todo models.Todo
			if err := cursor.Decode(&todo); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
			todos = append(todos, todo)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": todos})
	}
}
