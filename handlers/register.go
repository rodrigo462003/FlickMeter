package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rodrigo462003/FlickMeter/views/templates"
)

func RegisterGetHandler(c echo.Context) error {
	return Render(c, http.StatusOK, templates.Register())
}

func RegisterPostHandler(c echo.Context) error {
	return Render(c, http.StatusCreated, templates.BaseBody())
}

