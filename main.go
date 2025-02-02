package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rodrigo462003/FlickMeter/db"
	"github.com/rodrigo462003/FlickMeter/email"
	"github.com/rodrigo462003/FlickMeter/handlers"
	"github.com/rodrigo462003/FlickMeter/store"
	"golang.org/x/time/rate"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	connSTR := os.Getenv("CONN_STR")

	gmailPw := os.Getenv("GMAIL_APP_PW")
	from := os.Getenv("EMAIL")
	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	es := email.NewMailSender(from, gmailPw, host, port)
	d := db.New(connSTR)
	us := store.NewUserStore(d)
	h := handlers.NewHandler(us, es)

	e := echo.New()
	e.Debug = true
	e.Use(middleware.Secure())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(6))))
	e.Static("/public", "./public")
	e.GET("/", handlers.GetHome)

	e.GET("/signIn", h.UserHandler.GetSignIn)
	e.POST("/signIn", h.UserHandler.PostSignIn)
	e.GET("/register", h.UserHandler.GetRegister)
	e.POST("/register", h.UserHandler.PostRegister)
	e.POST("/register/username", h.UserHandler.PostUsername)
	e.POST("/register/email", h.UserHandler.PostEmail)
	e.POST("/register/password", h.UserHandler.PostPassword)
	e.POST("/register/verify", h.UserHandler.PostVerify)

	e.Logger.Fatal(e.Start(os.Getenv("ADDR")))
}
