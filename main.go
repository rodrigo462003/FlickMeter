package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rodrigo462003/FlickMeter/db"
	"github.com/rodrigo462003/FlickMeter/handlers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	connSTR := os.Getenv("CONN_STR")
	d := db.New(connSTR)
	h := handlers.NewHandler(d)

	e := echo.New()
	e.Debug = true
	e.Use(middleware.Secure())
	e.Static("/public", "./public")
	e.GET("/", handlers.GetHome)

	e.GET("/signIn", h.UserHandler.GetSignIn)
	e.POST("/signIn", h.UserHandler.PostSignIn)
	e.GET("/register", h.UserHandler.GetRegister)
	e.POST("/register", h.UserHandler.PostRegister)
	e.POST("/register/username", h.UserHandler.PostUsername)
	e.POST("/register/email", h.UserHandler.PostEmail)
	e.POST("/register/password", h.UserHandler.PostPassword)

	e.Logger.Fatal(e.Start(os.Getenv("ADDR")))
}
