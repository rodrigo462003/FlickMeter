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

	db := db.New(env["CONN_STR"])

	emailSender := email.NewMailSender(env["EMAIL_ADDR"], env["EMAIL_PW"], env["EMAIL_HOST"], env["EMAIL_PORT"])
	sessionStore := store.NewSessionStore(env["REDIS_ADDR"], db)
	userStore := store.NewUserStore(db)
	userService := service.NewUserService(userStore, sessionStore, emailSender)
	userHandler := handlers.NewUserHandler(userService)

	movieService := service.NewMovieService(env["API_KEY"])
	movieHandler := handlers.NewMovieHandler(movieService)

	e := echo.New()
	e.Debug = true
	limiterStore := middleware.NewRateLimiterMemoryStore(rate.Limit(10))
	e.Use(middleware.Secure(),
		middleware.Recover(),
		middleware.RateLimiter(limiterStore))

	e.Static("/public", "./public")
	e.GET("/", handlers.GetHome, userHandler.AuthMiddleware())

	userHandler.Register(e.Group("/user"))

	movieHandler.Register(e.Group("/movie"))

	e.Logger.Fatal(e.Start(env["ADDR"]))
}
