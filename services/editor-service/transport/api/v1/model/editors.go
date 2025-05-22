package model

type Editor struct {
	Email     string `json:"email" validate:"email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password"`
}
