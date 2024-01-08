package controllers

// "go.mongodb.org/mongo-driver/bson"
// "go.mongodb.org/mongo-driver/bson/primitive"
// "go.mongodb.org/mongo-driver/mongo"

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func generateToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func createSendToken(c *gin.Context, id string) error {
	token, err := generateToken(id)
	if err != nil {
		return err
	}

	c.SetCookie("jwt", token, int(time.Now().Add(time.Hour * 24 * 30).Unix()), "", "", true, true)
	return nil
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
			return
		}

		user.Password = hash
		user.ID = primitive.NewObjectID()

		inserted, err := collection.InsertOne(context.TODO(), user)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				c.JSON(http.StatusConflict, gin.H{"error": "User with this username is already exist"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := createSendToken(c, inserted.InsertedID.(primitive.ObjectID).Hex()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
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

		if err := createSendToken(c, existingUser.ID.Hex()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}
