package handlers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	uH "github.com/rodrigo462003/FlickMeter/handlers/user"
	"gorm.io/gorm"
)

type Handler struct {
	UserHandler uH.UserHandler
}

func NewHandler(d *gorm.DB) *Handler {
	return &Handler{
		UserHandler: *uH.NewUserHandler(d),
	}
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
