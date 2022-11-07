package validator

import (
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"matheus.com/vgs/internal/model"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateCreateUserRequest(c echo.Context) (model.CreateUserRequest, error) {
	req := model.CreateUserRequest{}
	if err := c.Bind(&req); err != nil {
		return req, err
	}
	return req, validate.Struct(&req)
}

func ValidateGetByIdRequest(c echo.Context) (uuid.UUID, error) {
	id := c.Param("id")
	return uuid.Parse(id)
}
