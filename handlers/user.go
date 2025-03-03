package handlers

import (
	"errors"
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

type signInForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
	Remember bool   `form:"remember"`
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
	form := &signInForm{}
	if err := c.Bind(form); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	userID, err := h.service.SignIn(form.Email, form.Password)
	if err != nil {
		if err, ok := err.(service.ValidationError); ok {
			return c.String(statusCode(err), err.Message())
		}
		return err
	}

	session, err := h.service.CreateSession(userID)
	if err != nil {
		return err
	}
	c.SetCookie(NewCookieSession(session))

	auth, err := h.service.CreateAuth(userID)
	if err != nil {
		return err
	}
	c.SetCookie(NewCookieAuth(auth))

	return c.NoContent(http.StatusOK)
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

	session, err := h.service.Verify(code, form.Username, form.Email, form.Password)
	if err != nil {
		switch e := err.(type) {
		case service.ValidationErrors:
			return Render(c, priorityStatusCode(e), templates.FormInvalid(e.FieldToMessage()))
		case service.ValidationError:
			return c.String(statusCode(e), e.Message())
		}
		return err
	}

	c.SetCookie(NewCookieSession(session))
	c.Response().Header().Set("HX-Redirect", "/")

	return c.NoContent(http.StatusCreated)
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

	if err, ok := err.(service.ValidationError); ok {
		return c.String(statusCode(err), err.Message())
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

func (h *userHandler) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("isAuth", false)
			cookie, err := c.Request().Cookie("session")
			if err == nil {
				if user, err := h.service.GetUserFromSession(cookie.Value); err == nil {
					c.Set("isAuth", true)
					c.Set("user", user)
					return next(c)
				}
			}

			auth, err := c.Request().Cookie("auth")
			if err != nil {
				return next(c)
			}

			user, err := h.service.GetUserFromAuth(auth.Value)
			if err != nil {
				return next(c)
			}

			if session, err := h.service.CreateSession(user.ID); err == nil {
				c.SetCookie(NewCookieSession(session))
			}

			c.Set("isAuth", true)
			c.Set("user", user)
			return next(c)
		}
	}
}

func (h *userHandler) AuthRequiredMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, err := c.Request().Cookie("session")
			if err == nil {
				if user, err := h.service.GetUserFromSession(session.Value); err != nil {
					c.Set("user", user)
					return next(c)
				}
			}

			auth, err := c.Request().Cookie("auth")
			if err != nil {
				return echo.ErrUnauthorized.WithInternal(err)
			}

			user, err := h.service.GetUserFromAuth(auth.Value)
			if err != nil {
				return echo.ErrUnauthorized.WithInternal(err)
			}

			if session, err := h.service.CreateSession(user.ID); err == nil {
				c.SetCookie(NewCookieSession(session))
			}

			c.Set("user", user)
			return next(c)
		}
	}
}
