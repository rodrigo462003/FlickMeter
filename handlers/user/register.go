package userHandler

import (
	"fmt"
	"net/http"
	"strconv"
	"unicode"

	"github.com/labstack/echo/v4"
	h "github.com/rodrigo462003/FlickMeter/handlers"
	"github.com/rodrigo462003/FlickMeter/views/templates"
)

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
	if valid, message := rForm.Username.isValid(); !valid {
		vm.Username = message
	}
	return true, vMessages{}
}

type username string

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
	return h.Render(c, http.StatusOK, templates.Register())
}

func (uH UserHandler)PostRegister(c echo.Context) error {
	var rForm registerForm
	err := c.Bind(&rForm)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if valid, messages := rForm.isValid(); !valid {
		return h.Render(c, http.StatusUnprocessableEntity, templates.Home(messages.Username == "2"))
	}
	fmt.Println(rForm)

	return h.Render(c, http.StatusCreated, templates.BaseBody())
}

func (uH UserHandler) PostUsername(c echo.Context) error {
	username := username(c.FormValue("username"))
	if valid, message := username.isValid(); !valid {
		return c.String(http.StatusUnprocessableEntity, message)
	}
	return c.String(http.StatusOK, "")
}
