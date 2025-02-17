package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rodrigo462003/FlickMeter/service"
	"github.com/rodrigo462003/FlickMeter/views/templates"
)

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

	err := h.service.VerifyUser(code, c.FormValue("email"))
	if err != nil {
		return c.NoContent(http.StatusOK)
	}

	if vErr, ok := err.(service.ValidationError); ok {
		return c.String(statusCode(vErr), vErr.Message())
	}

	return c.NoContent(http.StatusInternalServerError)
}

func (h userHandler) PostRegister(c echo.Context) error {
	form := struct {
		Username string `form:"username"`
		Email    string `form:"email"`
		Password string `form:"password"`
	}{}
	if err := c.Bind(&form); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := h.service.CreateUser(form.Username, form.Email, form.Password); err != nil {
		var sErr model.StatusErrors

		if errors.As(err, &sErr) {
			vm := sErr.Map()
			return Render(c, sErr.StatusCode(), templates.FormInvalid(vm))
		}
		return c.String(err.StatusCode(), err.Error())
	}

	return Render(c, http.StatusCreated, templates.FormVerifyCode(rForm.Email))
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
	err := h.service.ValidateEmail(c.FormValue("email"))
	if err == nil {
		return c.NoContent(http.StatusOK)
	}

	if vErr, ok := err.(service.ValidationError); ok {
		return c.String(statusCode(vErr), vErr.Message())
	}

	c.Logger().Error(err)
	return c.NoContent(http.StatusInternalServerError)
}

func (h *userHandler) PostPassword(c echo.Context) error {
	err := h.service.ValidatePassword(c.FormValue("password"))
	if err == nil {
		return c.NoContent(http.StatusOK)
	}

	if vErr, ok := err.(service.ValidationError); ok {
		return c.String(statusCode(vErr), vErr.Message())
	}

	c.Logger().Error(err)
	return c.NoContent(http.StatusInternalServerError)
}
