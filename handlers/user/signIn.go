package userHandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	h "github.com/rodrigo462003/FlickMeter/handlers"
	"github.com/rodrigo462003/FlickMeter/views/templates"
)

func (uH UserHandler) GetSignIn(c echo.Context) error {
	return h.Render(c, http.StatusOK, templates.SignIn())
}

func (uH UserHandler) PostSignIn(c echo.Context) error {
	return h.Render(c, http.StatusOK, templates.BaseBody())
}
