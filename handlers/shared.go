package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/rodrigo462003/FlickMeter/service"
	"github.com/rodrigo462003/FlickMeter/views/templates"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func statusCode(err service.ValidationError) int {
	switch {
	case err.Is(service.ErrConflict):
		return http.StatusConflict
	case err.Is(service.ErrUnprocessable):
		return http.StatusUnprocessableEntity
	default:
		panic("This shouldn't happen, All ValidationError types should be covered by the previous cases.")
	}
}

func GetHome(c echo.Context) error {
	return Render(c, http.StatusOK, templates.Home(false))
}
