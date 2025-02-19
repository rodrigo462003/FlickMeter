package handlers

import (
	"math"
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

func priorityStatusCode(err service.ValidationErrors) int {
	priorityMap := map[int]int{
		http.StatusConflict:            0,
		http.StatusUnprocessableEntity: 1,
	}

	prio, pCode := math.MaxInt, http.StatusUnprocessableEntity
	for _, err := range err.FieldToError() {
		sc := statusCode(err)
		if r, ok := priorityMap[sc]; ok {
			if r < prio {
				prio, pCode = r, sc
			}
		} else {
			panic("This shouldn't happen, All ValidationError types should be covered by priorityMap.")
		}
	}

	return pCode
}

func GetHome(c echo.Context) error {
	return Render(c, http.StatusOK, templates.Home(false))
}
