package controllers

// "go.mongodb.org/mongo-driver/bson"
// "go.mongodb.org/mongo-driver/bson/primitive"
// "go.mongodb.org/mongo-driver/mongo"

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/morf1lo/TODO-app/internal/models"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func verifyPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func CreateUser(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := user.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
			return
		}

		hash, err := hashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		user.Password = hash

		_, err = collection.InsertOne(context.TODO(), user, nil)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				c.JSON(http.StatusConflict, gin.H{"error": "User with this username is already exist"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	}
}

func Login(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var existingUser models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := collection.FindOne(context.TODO(), bson.D{{Key: "username", Value: user.Username}}).Decode(&existingUser)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusConflict, gin.H{"error": "Invalid username or password"})
				return
			}
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		matchPassword := verifyPassword(user.Password, existingUser.Password)
		if !matchPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Log in successfully"})
	}
}
