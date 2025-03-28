package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rodrigo462003/FlickMeter/service"
	"github.com/rodrigo462003/FlickMeter/views/templates"
)

type movieHandler struct {
	service service.MovieService
}

func NewMovieHandler(s service.MovieService) *movieHandler {
	return &movieHandler{s}
}

func (h *movieHandler) Register(g *echo.Group) {
	g.GET("/:id", h.Get)
	g.POST("/search", h.Search)
}

func (h *movieHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if len(id) < 0 {
		return echo.ErrBadRequest.WithInternal(errors.New("No movie ID param."))
	}

	movie, err := h.service.Get(id)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	return Render(c, http.StatusOK, templates.Movie(*movie, true))
}

func (h *movieHandler) Search(c echo.Context) error {
	query := c.FormValue("search")

	movies, err := h.service.Search(query)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	movies = movies[:min(5, len(movies))]
	return Render(c, http.StatusOK, templates.Results(movies))
}
