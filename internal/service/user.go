package service

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"matheus.com/vgs/api"
	"matheus.com/vgs/internal/model"
	"matheus.com/vgs/internal/repo"
)

type userService struct {
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) api.Users {
	return &userService{userRepo: userRepo}
}

func (svc *userService) Create(c echo.Context, user model.User) error {
	if err := svc.userRepo.Create(user); err != nil {
		return c.JSON(http.StatusInternalServerError, model.CreateUserResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, model.CreateUserResponse{Message: "user has been created"})
}

func (svc *userService) GetById(c echo.Context, id uuid.UUID) error {
	user, err := svc.userRepo.GetById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.GetByIdResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, model.GetByIdResponse{User: &user})
}
