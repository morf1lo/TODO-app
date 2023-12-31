package models

import "github.com/go-playground/validator/v10"

type User struct {
	Email			string	`json:"email" validate:"required,email"`
	Password	string	`json:"password" validate:"required,min=8"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
