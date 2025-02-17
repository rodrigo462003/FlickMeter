package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rodrigo462003/FlickMeter/db"
	"github.com/rodrigo462003/FlickMeter/email"
	"github.com/rodrigo462003/FlickMeter/handlers"
	"github.com/rodrigo462003/FlickMeter/service"
	"github.com/rodrigo462003/FlickMeter/store"
	"golang.org/x/time/rate"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	connSTR := os.Getenv("CONN_STR")
	gmailPw, from := os.Getenv("GMAIL_APP_PW"), os.Getenv("EMAIL")
	host, port := os.Getenv("EMAIL_HOST"), os.Getenv("EMAIL_PORT")

	es := email.NewMailSender(from, gmailPw, host, port)
	db := db.New(connSTR)
	us := store.NewUserStore(db)
	uh := handlers.NewUserHandler(service.NewUserService(us, es))

	e := echo.New()
	e.Debug = true
	e.Use(middleware.Secure())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(6))))
	e.Static("/public", "./public")
	e.GET("/", handlers.GetHome)

	e.GET("/signIn", uh.GetSignIn)
	e.POST("/signIn", uh.PostSignIn)
	e.GET("/register", uh.GetRegister)
	e.POST("/register", uh.PostRegister)
	e.POST("/register/username", uh.PostUsername)
	e.POST("/register/email", uh.PostEmail)
	e.POST("/register/password", uh.PostPassword)
	e.POST("/register/verify", uh.PostVerify)

	e.Logger.Fatal(e.Start(os.Getenv("ADDR")))
}
