package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rodrigo462003/FlickMeter/model"
	"github.com/rodrigo462003/FlickMeter/views/templates"
)

type UserHandler struct {
	userStore model.UserStore
}

func NewUserHandler(us model.UserStore) *UserHandler {
	return &UserHandler{userStore: us}
}

func (uH UserHandler) GetSignIn(c echo.Context) error {
	return Render(c, http.StatusOK, templates.SignIn())
}

func (uH UserHandler) PostSignIn(c echo.Context) error {
	return Render(c, http.StatusOK, templates.BaseBody())
}

type registerForm struct {
	Username string `form:"username"`
	Email    string `form:"email"`
	Password string `form:"password"`
	Confirm  string `form:"confirm"`
}

func (uH UserHandler) GetRegister(c echo.Context) error {
	return Render(c, http.StatusOK, templates.Register())
}

func (uH UserHandler) PostRegister(c echo.Context) error {
	var rForm registerForm
	err := c.Bind(&rForm)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	uErr := model.NewUser(rForm.Username, rForm.Email, rForm.Password, rForm.Confirm, uH.userStore)
	if uErr != nil {
		if uErr.StatusCode() == http.StatusInternalServerError {
			return c.NoContent(uErr.StatusCode())
		}
		vm := map[string]string{
			"username": uErr.Errors().Username,
			"email":    uErr.Errors().Email,
			"password": uErr.Errors().Password,
			"confirm":  uErr.Errors().Confirm,
		}
		Render(c, uErr.StatusCode(), templates.FormInvalid(vm))
	}

	return c.NoContent(http.StatusCreated)
}

func (uH *UserHandler) PostUsername(c echo.Context) error {
	err := model.ValidUsername(c.FormValue("username"), uH.userStore)
	if err != nil {
		if err.StatusCode() == http.StatusInternalServerError {
			c.NoContent(err.StatusCode())
		}
		return c.String(err.StatusCode(), err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (uH *UserHandler) PostEmail(c echo.Context) error {
	err := model.ValidEmail(c.FormValue("email"))
	if err != nil {
		return c.String(err.StatusCode(), err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (uH *UserHandler) PostPassword(c echo.Context) error {
	password := c.FormValue("password")
	confirm := c.FormValue("confirm")
	err, isDifferent := model.ValidPassword(password, confirm)
	if err != nil {
		if isDifferent {
			return Render(c, http.StatusUnprocessableEntity, templates.Oob("confirmErr", "* Passwords don't match."))
		}
		c.String(err.StatusCode(), err.Error())
		return Render(c, err.StatusCode(), templates.Oob("confirmErr", ""))
	}

	return Render(c, http.StatusOK, templates.Oob("confirmErr", ""))
}
