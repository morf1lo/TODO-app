package models

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID				primitive.ObjectID	`bson:"_id"`
	UserID		primitive.ObjectID	`json:"userId" validate:"required"`
	Title			string							`json:"title" validate:"required,min=1,max=30"`
	Text			string							`json:"text" validate:"required,min=1,max=120"`
	Completed	bool								`json:"completed"`
}

func (t *Todo) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}
