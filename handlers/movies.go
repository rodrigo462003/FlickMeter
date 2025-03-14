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
	g.GET("/:id", h.GetMovie)
    g.POST("/search", h.SearchMovies)
}

func (h *movieHandler) GetMovie(c echo.Context) error {
	id := c.Param("id")
	if len(id) < 0 {
		return echo.ErrBadRequest.WithInternal(errors.New("No movie ID param."))
	}

	movie, err := h.service.GetMovie(id)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	return Render(c, http.StatusOK, templates.Movie(*movie, true))
}

func (h *movieHandler) SearchMovies(c echo.Context) error {
	search := c.FormValue("search")

    movies := h.service.SearchMovies(search)

    return c.JSON(http.StatusOK, movies)
}
