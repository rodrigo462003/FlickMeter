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
	limiterStore := middleware.NewRateLimiterMemoryStore(rate.Limit(10))
	e.Use(middleware.Secure(),
		middleware.Recover(),
		middleware.RateLimiter(limiterStore))

	e.Static("/public", "./public")
	e.GET("/", handlers.GetHome)

	userHandler.Register(e.Group("/user"))

	e.Logger.Fatal(e.Start(env["ADDR"]))
}
