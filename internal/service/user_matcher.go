package service

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"matheus.com/vgs/api"
	"matheus.com/vgs/internal/engine"
	"matheus.com/vgs/internal/logger"
	"matheus.com/vgs/internal/messaging"
	"matheus.com/vgs/internal/repo"
)

type userMatcher struct {
	userRepo        repo.UserRepo
	sqsPublisher    *messaging.SQSPublisher
	evaluatorEngine *engine.EvaluatorEngine
}

func NewUserMatcherService(userRepo repo.UserRepo, sqsPublisher *messaging.SQSPublisher, engine *engine.EvaluatorEngine) api.UserMatcher {
	matcher := &userMatcher{
		userRepo:        userRepo,
		sqsPublisher:    sqsPublisher,
		evaluatorEngine: engine,
	}
	return matcher
}

func (svc *userMatcher) Match(c echo.Context, filter map[string]interface{}) error {
	go func() {
		users, err := svc.userRepo.GetAll()
		if err != nil {
			logger.Logger().Error("error ", err)
			return
		}
		filter["Users"] = users
		users, err = svc.evaluatorEngine.Evaluate(filter)
		if err != nil {
			logger.Logger().Error("could not match users: ", err)
		}

		for _, user := range users {
			svc.sqsPublisher.Publish(user)
		}
	}()
	return c.JSON(http.StatusAccepted, "")
}
