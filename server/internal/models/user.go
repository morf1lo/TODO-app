package models

import "github.com/go-playground/validator/v10"

type User struct {
	Username	string	`json:"username" validate:"required,min=3,max=16"`
	Password	string	`json:"password" validate:"required,min=8"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
