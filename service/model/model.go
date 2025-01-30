package model

type UserCreate struct {
	Username string
	Email    string
	Password string
	Role     string
}
