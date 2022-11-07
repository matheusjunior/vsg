package api

import "github.com/labstack/echo/v4"

type UserMatcher interface {
	Match(c echo.Context, filter map[string]interface{}) error
}
