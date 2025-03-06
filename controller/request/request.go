package request

import "github.com/denizdoganinsider/kpi_project/service/model"

type AddUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (addUserRequest AddUserRequest) ToModel() model.UserCreate {
	return model.UserCreate{
		Username: addUserRequest.Username,
		Email:    addUserRequest.Email,
		Password: addUserRequest.Password,
		Role:     addUserRequest.Role,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
