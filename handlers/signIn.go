package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rodrigo462003/FlickMeter/views/templates"
)

func SignInGetHandler(c echo.Context) error {
	return Render(c, http.StatusOK, templates.SignIn())
}

func SignInPostHandler(c echo.Context) error {
	return Render(c, http.StatusOK, templates.BaseBody())
}
