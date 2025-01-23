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

type registerForm struct {
	Username username `form:"username"`
	Email    email    `form:"email"`
	Password password `form:"password"`
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

type username string

func (u username) isValid() (bool, string) {
	const maxLen = 15
	if len(u) > maxLen {
		return false, fmt.Sprint("* Username must have at most ", strconv.Itoa(maxLen), " characters.")
	}
	if len(u) == 0 {
		return false, "* Username is required."
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
		return false, fmt.Sprint("* ", string(r), " is not allowed.")
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
		return c.String(http.StatusConflict, "* Username already exists.")
	}

	return c.NoContent(http.StatusOK)
}

type email string

func (e email) isValid() (bool, string) {
	if len(e) == 0 {
		return false, "* Email address is required."
	}

	if _, err := mail.ParseAddress(string(e)); err != nil {
		return false, "* This is not a valid email address."
	}

	return true, ""
}

func (uH *UserHandler) PostEmail(c echo.Context) error {
	emailS := email(c.FormValue("email"))
	if valid, message := emailS.isValid(); !valid {
		return c.String(http.StatusUnprocessableEntity, message)
	}

	alreadyExists, err := uH.us.EmailExists(string(emailS))
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		panic(err)
	}
	if alreadyExists {
		return c.String(http.StatusConflict, "* This email is already registered.")
	}

	return c.NoContent(http.StatusOK)
}

type password string

func (p password) isValid() (bool, string) {
	n := len(p)
	if n < 8 {
		return false, "* Password must contain atleast 8 characters."
	}

	if n > 64 {
		return false, "* Password must contain at most 64 characters."
	}

	const MinPrintableASCII = 32
	const MaxPrintableASCII = unicode.MaxASCII - 1
	for _, c := range p {
		if c > MaxPrintableASCII || c < MinPrintableASCII {
			return false, fmt.Sprintf("* Your password must not contain: '%c'. Use only letters, numbers, and common symbols like !, @, #, *, (, ), etc.", c)
		}
	}

	return true, ""
}

func (uH *UserHandler) PostPassword(c echo.Context) error {
	password := password(c.FormValue("password"))
	if valid, message := password.isValid(); !valid {
		c.String(http.StatusUnprocessableEntity, message)
		return Render(c, http.StatusUnprocessableEntity, templates.Oob("confirmErr", ""))
	}

	confirm := c.FormValue("confirm")
	if string(password) != confirm && len(confirm) > 0 {
		return Render(c, http.StatusConflict, templates.Oob("confirmErr", "* Passwords don't match."))
	}

	return Render(c, http.StatusOK, templates.Oob("confirmErr", ""))
}
