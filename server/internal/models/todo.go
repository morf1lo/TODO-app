package models

import "github.com/go-playground/validator/v10"

type Todo struct {
	Username	string	`json:"username" validate:"required,min=3,max=16"`
	Title			string	`json:"title" validate:"required,min=1,max=30"`
	Text			string	`json:"text" validate:"required,min=1,max=120"`
	Completed	bool		`json:"completed"`
}

func (t *Todo) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}
