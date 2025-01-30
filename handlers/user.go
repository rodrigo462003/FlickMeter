package handlers

import (
	"net/http"

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
	return Render(c, http.StatusOK, templates.FormVerifyCode("gsdagsdagsdafkljasdfjs@jglskdj"))
}

func (uh UserHandler) PostRegister(c echo.Context) error {
	var rForm registerForm
	err := c.Bind(&rForm)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	uErr := model.NewUser(rForm.Username, rForm.Email, rForm.Password, uh.userStore, uh.emailSender)
	if uErr != nil {
		if uErr.StatusCode() == http.StatusInternalServerError {
			return c.NoContent(uErr.StatusCode())
		}
		vm := map[string]string{
			"username": uErr.Errors().Username,
			"email":    uErr.Errors().Email,
			"password": uErr.Errors().Password,
		}
		return Render(c, uErr.StatusCode(), templates.FormInvalid(vm))
	}

	return Render(c, http.StatusCreated, templates.FormVerifyCode(rForm.Email))
}

func (uh *UserHandler) PostUsername(c echo.Context) error {
	err := model.ValidUsername(c.FormValue("username"), uh.userStore)
	if err != nil {
		if err.StatusCode() == http.StatusInternalServerError {
			c.NoContent(err.StatusCode())
		}
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
