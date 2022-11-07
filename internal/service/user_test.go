package service

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserIt(t *testing.T) {
	client := resty.New()
	defaultRequest := func() *resty.Request {
		return client.R().
			SetHeader("Accept", "application/json").
			SetHeader("Content-Type", "application/json")
	}

	t.Run("user creation fails when email is empty or invalid", func(t *testing.T) {
		user := `
		{
			"name": "user",
			"email": "",
			"city": "campinas"
		}`
		resp, err := defaultRequest().
			SetBody(user).
			Post("http://localhost:8080/users")
		assert.Nil(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusBadRequest)

		user = `
		{
			"name": "user",
			"email": "invalid_email@",
			"city": "campinas"
		}`
		resp, err = defaultRequest().
			SetBody(user).
			Post("http://localhost:8080/users")
		assert.Nil(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusBadRequest)
	})

	t.Run("user creation fails when email is duplicated", func(t *testing.T) {
		user := `
		{
			"name": "user",
			"email": "user1@gmail.com",
			"city": "campinas"
		}`
		resp, err := defaultRequest().
			SetBody(user).
			Post("http://localhost:8080/users")
		assert.Nil(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusInternalServerError)
	})
}

func TestGetUserByIdIt(t *testing.T) {
	client := resty.New()
	defaultRequest := func() *resty.Request {
		return client.R().
			SetHeader("Accept", "application/json").
			SetHeader("Content-Type", "application/json")
	}

	t.Run("get user by fails with invalid UUIDv4", func(t *testing.T) {
		resp, err := defaultRequest().
			SetPathParam("id", "invalid").
			Get("http://localhost:8080/users/{id}")
		assert.Nil(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusBadRequest)
	})
}
