package main

import (
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
	env, err := godotenv.Read()
	if err != nil {
		panic(err)
	}

	sessionStore := store.NewSessionStore(env["REDIS_ADDR"])
	userStore := store.NewUserStore(db.New(env["CONN_STR"]))
	emailSender := email.NewMailSender(env["EMAIL_ADDR"], env["EMAIL_PW"], env["EMAIL_HOST"], env["EMAIL_PORT"])
	userService := service.NewUserService(userStore, sessionStore, emailSender)
	userHandler := handlers.NewUserHandler(userService)

	e := echo.New()
	e.Debug = true
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(10))))

	e.Static("/public", "./public")
	e.GET("/", handlers.GetHome)

	e.GET("/signIn", userHandler.GetSignIn)
	e.POST("/signIn", userHandler.PostSignIn)
	e.GET("/register", userHandler.GetRegister)
	e.POST("/register", userHandler.PostRegister)
	e.POST("/register/username", userHandler.PostUsername)
	e.POST("/register/email", userHandler.PostEmail)
	e.POST("/register/password", userHandler.PostPassword)
	e.POST("/register/verify", userHandler.PostVerify)

	e.Logger.Fatal(e.Start(env["ADDR"]))
}
