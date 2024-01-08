package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/morf1lo/TODO-app/internal/models"
)

func CreateTodo(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("uid")
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not authorized"})
			return
		}

		var todo models.Todo
		
		objectId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		todo.UserID = objectId

		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := todo.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo.ID = primitive.NewObjectID()

		inserted, err := collection.InsertOne(context.TODO(), todo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"todoID": inserted.InsertedID}})
	}
}

func FindAllTodos(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("uid")
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not authorized"})
			return
		}

		objectId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cursor, err := collection.Find(context.TODO(), bson.M{"userid": objectId})
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

func UpdateTodo(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("uid")
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not authorized"})
			return
		}

		var updateTodoOptions map[string]interface{}

		if err := c.ShouldBindJSON(&updateTodoOptions); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		objectId, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updatedTodo, err := collection.UpdateByID(context.TODO(), objectId, bson.D{{Key: "$set", Value: updateTodoOptions}})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": updatedTodo})
	}
}

func DeleteTodo(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("uid")
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not authorized"})
			return
		}

		userObjectId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todoId := c.Param("id")

		todoObjectId, err := primitive.ObjectIDFromHex(todoId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		deletedResult, err := collection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: todoObjectId}, {Key: "userid", Value: userObjectId}})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": deletedResult})
	}
}
