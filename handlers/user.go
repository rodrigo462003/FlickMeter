package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rodrigo462003/FlickMeter/service"
	"github.com/rodrigo462003/FlickMeter/views/templates"
)

type registerForm struct {
	Username string `form:"username"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

type userHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *userHandler {
	return &userHandler{s}
}

func (h *userHandler) Register(g *echo.Group) {
	g.GET("/signIn", h.GetSignIn)
	g.POST("/signIn", h.PostSignIn)
	g.GET("/register", h.GetRegister)
	g.POST("/register", h.PostRegister)
	g.POST("/register/username", h.PostUsername)
	g.POST("/register/email", h.PostEmail)
	g.POST("/register/password", h.PostPassword)
	g.POST("/register/verify", h.PostVerify)
}

func (h *userHandler) GetSignIn(c echo.Context) error {
	return Render(c, http.StatusOK, templates.SignIn())
}

func (h *userHandler) PostSignIn(c echo.Context) error {
	return Render(c, http.StatusOK, templates.BaseBody())
}

func (h *userHandler) GetRegister(c echo.Context) error {
	return Render(c, http.StatusOK, templates.Register())
}

func (h *userHandler) PostVerify(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	digits := c.Request().Form["code"]
	if len(digits) != 6 {
		return echo.ErrBadRequest.WithInternal(errors.New("Form needs to have 6 code fields."))
	}
	code := strings.Join(digits, "")
	if len(code) != 6 {
		return c.NoContent(http.StatusUnprocessableEntity)
	}

	form := &registerForm{}
	if err := c.Bind(form); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	uuid, err := h.service.Verify(code, form.Username, form.Email, form.Password)
	if err != nil {
		switch e := err.(type) {
		case service.ValidationErrors:
			return Render(c, priorityStatusCode(e), templates.FormInvalid(e.FieldToMessage()))
		case service.ValidationError:
			return c.String(statusCode(e), e.Message())
		}
		return err
	}

	c.SetCookie(newCookie(uuid))
	c.Response().Header().Set("HX-Redirect", "/")

	return c.NoContent(http.StatusCreated)
}

func newCookie(uuid uuid.UUID) *http.Cookie {
	return &http.Cookie{
		Name:     "session",
		Value:    uuid.String(),
		Path:     "/",
		Secure:   false, //SET TO SECURE FOR HTTPS.
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func (h *userHandler) PostRegister(c echo.Context) error {
	form := &registerForm{}
	if err := c.Bind(form); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := h.service.Register(form.Username, form.Email, form.Password); err != nil {
		switch e := err.(type) {
		case service.ValidationErrors:
			return Render(c, priorityStatusCode(e), templates.FormInvalid(e.FieldToMessage()))
		case service.ValidationError:
			return Render(c, statusCode(e), templates.FormVerifyCode(form.Email, e.Message()))
		}
		return err
	}

	return Render(c, http.StatusCreated, templates.FormVerifyCode(form.Email, ""))
}

func (h *userHandler) PostUsername(c echo.Context) error {
	err := h.service.ValidateUsername(c.FormValue("username"))
	if err == nil {
		return c.NoContent(http.StatusOK)
	}

	if vErr, ok := err.(service.ValidationError); ok {
		return c.String(statusCode(vErr), vErr.Message())
	}

	return err
}

func (h *userHandler) PostEmail(c echo.Context) error {
	if err := h.service.ValidateEmail(c.FormValue("email")); err != nil {
		return c.String(statusCode(err), err.Message())
	}
	return c.NoContent(http.StatusOK)
}

func (h *userHandler) PostPassword(c echo.Context) error {
	if err := h.service.ValidatePassword(c.FormValue("password")); err != nil {
		return c.String(statusCode(err), err.Message())
	}

	return c.NoContent(http.StatusOK)
}
