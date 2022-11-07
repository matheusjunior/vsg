package model

import "time"

func CreateUserRequestToUserModel(req CreateUserRequest) User {
	return User{
		Name:      req.Name,
		Email:     req.Email,
		City:      req.City,
		LastLogin: time.Unix(req.LastLogin, 0),
	}
}
