package handlers

import (
	"net/http"
	"strings"

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

func NewUserHandler(us service.UserService) *userHandler {
	return &userHandler{us}
}

func (h userHandler) GetSignIn(c echo.Context) error {
	return Render(c, http.StatusOK, templates.SignIn())
}

func (h userHandler) PostSignIn(c echo.Context) error {
	return Render(c, http.StatusOK, templates.BaseBody())
}

func (h userHandler) GetRegister(c echo.Context) error {
	return Render(c, http.StatusOK, templates.Register())
}

func (h *userHandler) PostVerify(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	digits := c.Request().Form["code"]
	if len(digits) != 6 {
		return c.NoContent(http.StatusBadRequest)
	}
	code := strings.Join(digits, "")
	if len(code) != 6 {
		return c.NoContent(http.StatusUnprocessableEntity)
	}

	form := &registerForm{}
	if err := c.Bind(form); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	cookie, err := h.service.Verify(code, form.Username, form.Email, form.Password)
	if err != nil {
		switch e := err.(type) {
		case service.ValidationErrors:
			return Render(c, priorityStatusCode(e), templates.FormInvalid(e.FieldToMessage()))
		case service.ValidationError:
			return c.String(statusCode(e), e.Message())
		default:
			c.Logger().Error(err)
			c.NoContent(http.StatusInternalServerError)
		}
	}

	c.SetCookie(cookie)
	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusCreated)
}

func (h userHandler) PostRegister(c echo.Context) error {
	form := &registerForm{}
	if err := c.Bind(form); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := h.service.Register(form.Username, form.Email, form.Password); err != nil {
		switch e := err.(type) {
		case service.ValidationErrors:
			return Render(c, priorityStatusCode(e), templates.FormInvalid(e.FieldToMessage()))
		case service.ValidationError:
			return Render(c, statusCode(e), templates.FormVerifyCode(form.Email, e.Message()))
		default:
			c.Logger().Error(err)
			return c.NoContent(http.StatusInternalServerError)
		}
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

	c.Logger().Error(err)
	return c.NoContent(http.StatusInternalServerError)
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
