package model

type UserCreate struct {
	Username     string
	Email        string
	PasswordHash string
	Role         string
}
