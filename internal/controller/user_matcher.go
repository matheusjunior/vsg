package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"matheus.com/vgs/api"
	"matheus.com/vgs/internal/model"
)

type UserMatcherController struct {
	userMatcherSvc api.UserMatcher
}

func NewUserMatcherController(userMatcherSvc api.UserMatcher) *UserMatcherController {
	return &UserMatcherController{userMatcherSvc: userMatcherSvc}
}

func (controller *UserMatcherController) Match(c echo.Context) error {
	var filter map[string]interface{}
	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, model.MatchResponse{Error: err.Error()})
	}
	if filter == nil {
		return c.JSON(http.StatusBadRequest, model.MatchResponse{Message: "missing filter payload"})
	}
	return controller.userMatcherSvc.Match(c, filter)
}
