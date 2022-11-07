package model

const MsgUserCreated = "user has been created"

type CreateUserRequest struct {
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	City      string `json:"city"`
	LastLogin int64  `json:"last_login"`
}
