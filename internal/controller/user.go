package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"matheus.com/vgs/api"
	"matheus.com/vgs/internal/model"
	"matheus.com/vgs/internal/validator"
)

type UserController struct {
	userSvc api.Users
}

func NewUserController(userSvc api.Users) *UserController {
	return &UserController{userSvc: userSvc}
}

func (controller *UserController) Create(c echo.Context) error {
	req, err := validator.ValidateCreateUserRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.CreateUserResponse{Error: err.Error()})
	}
	user := model.CreateUserRequestToUserModel(req)
	return controller.userSvc.Create(c, user)
}

func (controller *UserController) GetById(c echo.Context) error {
	id, err := validator.ValidateGetByIdRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.GetByIdResponse{Error: err.Error()})
	}
	return controller.userSvc.GetById(c, id)
}
