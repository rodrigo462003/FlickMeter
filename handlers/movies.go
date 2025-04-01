package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rodrigo462003/FlickMeter/model"
	"github.com/rodrigo462003/FlickMeter/service"
	"github.com/rodrigo462003/FlickMeter/views/templates"
)

type movieHandler struct {
	service service.MovieService
}

func NewMovieHandler(s service.MovieService) *movieHandler {
	return &movieHandler{s}
}

func (h *movieHandler) Register(g *echo.Group, authMiddleware echo.MiddlewareFunc) {
	g.GET("/:id", h.Get, authMiddleware)
	g.POST("/search", h.Search)
	g.POST("/newReview", h.NewReview, authMiddleware)
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

	isAuth := c.Get("isAuth").(bool)
	return Render(c, http.StatusOK, templates.Movie(*movie, isAuth))
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

func (h *movieHandler) NewReview(c echo.Context) error {
	form := struct {
		Title string `form:"title"`
		Text  string `form:"text"`
		Movie uint   `form:"movieID"`
	}{}

	if err := c.Bind(&form); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	user, ok := c.Get("user").(*model.User)
	if !ok {
		return echo.ErrUnauthorized.WithInternal(errors.New("Failed to assert context 'user' as a *model.User"))
	}

	review, err := h.service.CreateReview(form.Title, form.Text, form.Movie, user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, "You already have a review for this movie.")
	}

	return Render(c, http.StatusCreated, templates.Review(review))
}
