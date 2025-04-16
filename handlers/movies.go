package handlers

import (
	"errors"
	"net/http"
	"slices"
	"strconv"

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
	g.GET("/:id/", h.Get, authMiddleware)
	g.POST("/search", h.Search)
	g.POST("/:id/review", h.UpdateReview, authMiddleware)
	g.GET("/:id/review", h.GetReview, authMiddleware)
}

func (h *movieHandler) Home(c echo.Context) error {
	user, ok := c.Get("user").(*model.User)
	if !ok {
		user = &model.User{}
	}

	topMovies := h.service.Top()

	return Render(c, http.StatusOK, templates.Home(topMovies, user))
}

func (h *movieHandler) GetReview(c echo.Context) error {
	user, ok := c.Get("user").(*model.User)
	if !ok {
		return echo.ErrUnauthorized.WithInternal(errors.New("Failed to assert context 'user' as a *model.User"))
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	review, err := h.service.GetReview(uint(id), user.ID)
	if err != nil {
		return err
	}

	return Render(c, http.StatusOK, templates.NewForm(review))
}

func (h *movieHandler) Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	movie, err := h.service.Get(uint(id))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	user, ok := c.Get("user").(*model.User)
	if !ok {
		user = &model.User{}
	}
	var userReview *model.Review
	if user.ID != 0 {
		i := slices.IndexFunc(movie.Reviews, func(r model.Review) bool {
			return r.UserID == user.ID
		})
		if i > -1 {
			movie.Reviews[i], movie.Reviews[len(movie.Reviews)-1] = movie.Reviews[len(movie.Reviews)-1], movie.Reviews[i]
			userReview = &movie.Reviews[len(movie.Reviews)-1]
			movie.Reviews = movie.Reviews[:len(movie.Reviews)-1]
		}
	}
	return Render(c, http.StatusOK, templates.Movie(*movie, user, userReview))
}

func (h *movieHandler) Search(c echo.Context) error {
	query := c.FormValue("search")

	movies, err := h.service.Search(query)
	if err != nil {
		return err
	}

	movies = movies[:min(5, len(movies))]
	return Render(c, http.StatusOK, templates.Results(movies))
}

func (h *movieHandler) UpdateReview(c echo.Context) error {
	form := struct {
		Title  string `form:"title"`
		Text   string `form:"text"`
		Rating uint   `form:"rating"`
	}{}
	if err := c.Bind(&form); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	user, ok := c.Get("user").(*model.User)
	if !ok {
		return echo.ErrUnauthorized.WithInternal(errors.New("Failed to assert context 'user' as a *model.User"))
	}

	review, err := h.service.UpdateReview(form.Title, form.Text, form.Rating, uint(id), user.ID)
	if err != nil {
		if err, ok := err.(service.ValidationError); ok {
			return c.String(statusCode(err), err.Message())
		}
		return err
	}

	return Render(c, http.StatusCreated, templates.Review(review))
}
