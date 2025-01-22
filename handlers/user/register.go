package userHandler

import (
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"unicode"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
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

type username string

type registerForm struct {
	Username username `form:"username"`
	Email    string   `form:"email"`
	Password string   `form:"password"`
	Confirm  string   `form:"confirm"`
}

type vMessages struct {
	Username string
	Email    string
	Password string
	Confirm  string
}

func (rForm registerForm) isValid() (bool, vMessages) {
	vm := vMessages{}
	valid, message := rForm.Username.isValid()
	if !valid {
		vm.Username = message
	}
	return true, vMessages{}
}

func (u username) isValid() (bool, string) {
	const maxLen = 15
	if len(u) > maxLen {
		return false, fmt.Sprint("Username must have at most ", strconv.Itoa(maxLen), " characters.")
	}
	if len(u) < 1 {
		return false, "Username is required."
	}
	for _, r := range u {
		if unicode.IsDigit(r) {
			continue
		}
		if unicode.IsLetter(r) {
			continue
		}
		if r == '_' {
			continue
		}
		return false, fmt.Sprint(string(r), " is not allowed.")
	}

	return true, ""
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

	if valid, messages := rForm.isValid(); !valid {
		return Render(c, http.StatusUnprocessableEntity, templates.Home(messages.Username == "2"))
	}
	fmt.Println(rForm)

	return Render(c, http.StatusCreated, templates.BaseBody())
}

func (uH *UserHandler) PostUsername(c echo.Context) error {
	username := username(c.FormValue("username"))
	if valid, message := username.isValid(); !valid {
		return c.String(http.StatusUnprocessableEntity, message)
	}

	alreadyExists, err := uH.us.UserNameExists(string(username))
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		panic(err)
	}
	if alreadyExists {
		return c.String(http.StatusConflict, "Username already exists.")
	}

	return c.String(http.StatusOK, "")
}

func (uH *UserHandler) PostEmail(c echo.Context) error {
	emailS := c.FormValue("email")
	if len(emailS) < 1 {
		return c.String(http.StatusUnprocessableEntity, "Email address is required.")
	}

	email, err := mail.ParseAddress(emailS)
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, "This is not a valid email address.")
	}

	emailS = email.Address
	alreadyExists, err := uH.us.EmailExists(emailS)
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		panic(err)
	}
	if alreadyExists {
		return c.String(http.StatusConflict, "Email already exists.")
	}

	return c.String(http.StatusOK, "")
}
