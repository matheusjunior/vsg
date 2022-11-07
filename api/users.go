package api

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"matheus.com/vgs/internal/model"
)

type Users interface {
	Create(c echo.Context, user model.User) error
	GetById(c echo.Context, id uuid.UUID) error
}
