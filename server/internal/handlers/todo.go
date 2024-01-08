package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/morf1lo/TODO-app/internal/models"
)

type UpdateOptions struct {
	Title     string `json:"title"`
	Text      string `json:"text"`
	Completed *bool   `json:"completed,omitempty"`
}

func (u *UpdateOptions) FilterUpdateOptions() bson.M {
	update := bson.M{"$set": bson.M{}}

	if strings.TrimSpace(u.Title) != "" {
		update["$set"].(bson.M)["title"] = strings.TrimSpace(u.Title)
	}

	if strings.TrimSpace(u.Text) != "" {
		update["$set"].(bson.M)["text"] = strings.TrimSpace(u.Text)
	}

	if u.Completed != nil {
		update["$set"].(bson.M)["completed"] = *u.Completed
	}

	return update
}

// Create Todo
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

// Get all Todos
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

// Update Todo info
func UpdateTodo(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("uid")
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not authorized"})
			return
		}

		var updateFields UpdateOptions

		if err := c.ShouldBindJSON(&updateFields); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		todoObjectId, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userObjectId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not authorized"})
		}

		update := updateFields.FilterUpdateOptions()

		data, err := collection.UpdateOne(context.TODO(), bson.M{"_id": todoObjectId, "userid": userObjectId}, update)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
	}
}

// Delete Todo
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
