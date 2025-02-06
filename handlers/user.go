package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rodrigo462003/FlickMeter/email"
	"github.com/rodrigo462003/FlickMeter/model"
	"github.com/rodrigo462003/FlickMeter/views/templates"
)

type UserHandler struct {
	userStore   model.UserStore
	emailSender email.EmailSender
}

func NewUserHandler(us model.UserStore, es email.EmailSender) *UserHandler {
	return &UserHandler{us, es}
}

func (uh UserHandler) GetSignIn(c echo.Context) error {
	return Render(c, http.StatusOK, templates.SignIn())
}

func (uh UserHandler) PostSignIn(c echo.Context) error {
	return Render(c, http.StatusOK, templates.BaseBody())
}

type registerForm struct {
	Username string `form:"username"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (uh UserHandler) GetRegister(c echo.Context) error {
	return Render(c, http.StatusOK, templates.Register())
}

func (uh UserHandler) PostRegister(c echo.Context) error {
	var rForm registerForm
	err := c.Bind(&rForm)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := model.NewUser(rForm.Username, rForm.Email, rForm.Password, uh.userStore, uh.emailSender); err != nil {
		var sErr model.StatusErrors
		if errors.As(err, &sErr) {
			vm := sErr.Map()
			return Render(c, sErr.StatusCode(), templates.FormInvalid(vm))
		}
		return c.String(err.StatusCode(), err.Error())
	}

	return Render(c, http.StatusCreated, templates.FormVerifyCode(rForm.Email))
}

func (uh *UserHandler) PostVerify(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	code := c.Request().Form["code"]
	if len(code) != 6 {
		return c.NoContent(http.StatusBadRequest)
	}
	codeString := strings.Join(code, "")
	if len(codeString) != 6 {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)
}

func (uh *UserHandler) PostUsername(c echo.Context) error {
	err := model.ValidUsername(c.FormValue("username"), uh.userStore)
	if err != nil {
		return c.String(err.StatusCode(), err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (uh *UserHandler) PostEmail(c echo.Context) error {
	err := model.ValidEmail(c.FormValue("email"))
	if err != nil {
		return c.String(err.StatusCode(), err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (uh *UserHandler) PostPassword(c echo.Context) error {
	password := c.FormValue("password")
	err := model.ValidPassword(password)
	if err != nil {
		return c.String(err.StatusCode(), err.Error())
	}

	return c.String(http.StatusOK, "")
}
