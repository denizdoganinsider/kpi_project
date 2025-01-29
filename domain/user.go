package domain

type User struct {
	Id           int64
	Username     string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    string
	UpdatedAt    string
}
