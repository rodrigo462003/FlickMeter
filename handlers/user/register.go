package userHandler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/rodrigo462003/FlickMeter/model"
	"github.com/rodrigo462003/FlickMeter/store"
	"github.com/rodrigo462003/FlickMeter/views/templates"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

type registerForm struct {
	Username string `form:"username"`
	Email    string `form:"email"`
	Password string `form:"password"`
	Confirm  string `form:"confirm"`
}

func (rForm registerForm) isValid(us store.UserStore) (map[string]string, int) {
	vm := make(map[string]string)
	statusCode := http.StatusOK

	vUsername, err := model.NewUsername(rForm.Username)
	if err != nil {
		vm["username"] = err.Error()
		statusCode = http.StatusUnprocessableEntity
	} else {
		alreadyExists, err := us.UserNameExists(vUsername)
		if err != nil {
			return vm, http.StatusInternalServerError
		}
		if alreadyExists {
			vm["username"] = "* Username already taken."
			if statusCode != http.StatusUnprocessableEntity {
				statusCode = http.StatusConflict
			}
		}
	}

	vEmail, err := model.NewEmail(rForm.Email)
	if err != nil {
		vm["email"] = err.Error()
		statusCode = http.StatusUnprocessableEntity
	}

	password := rForm.Password
	vPassword, err := model.NewPassword(password)
	if err != nil {
		vm["password"] = err.Error()
		statusCode = http.StatusUnprocessableEntity
	}

	confirm := rForm.Confirm
	if password != confirm && len(confirm) > 0 {
		vm["confirm"] = "* Passwords don't match."
		statusCode = http.StatusUnprocessableEntity
	}

	_ = model.NewUser(vUsername, vEmail, vPassword)

	return vm, statusCode
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

	vm, statusCode := rForm.isValid(*uH.us)
	if statusCode == http.StatusInternalServerError {
		return c.NoContent(http.StatusInternalServerError)
	}

	if statusCode != http.StatusOK {
		Render(c, statusCode, templates.FormInvalid(vm))
	}

	return c.NoContent(http.StatusCreated)
}

func (uH *UserHandler) PostUsername(c echo.Context) error {
	username, err := model.NewUsername(c.FormValue("username"))
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	alreadyExists, err := uH.us.UserNameExists(username)
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		c.Logger().Error(err)
	}
	if alreadyExists {
		return c.String(http.StatusConflict, "* Username already taken.")
	}

	return c.NoContent(http.StatusOK)
}

func (uH *UserHandler) PostEmail(c echo.Context) error {
	_, err := model.NewEmail(c.FormValue("email"))
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (uH *UserHandler) PostPassword(c echo.Context) error {
	password := c.FormValue("password")
	_, err := model.NewPassword(password)
	if err != nil {
		c.String(http.StatusUnprocessableEntity, err.Error())
		return Render(c, http.StatusUnprocessableEntity, templates.Oob("confirmErr", ""))
	}

	confirm := c.FormValue("confirm")
	if password != confirm && len(confirm) > 0 {
		return Render(c, http.StatusUnprocessableEntity, templates.Oob("confirmErr", "* Passwords don't match."))
	}

	return Render(c, http.StatusOK, templates.Oob("confirmErr", ""))
}
