package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/rodrigo462003/FlickMeter/handlers"
	uH "github.com/rodrigo462003/FlickMeter/handlers/user"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Static("/public", "./public")
	e.GET("/", handlers.HomeHandler)

	userH := uH.UserHandler{}
	e.GET("/signIn", userH.GetSignIn)
	e.POST("/signIn", userH.PostSignIn)
	e.GET("/register", userH.GetRegister)
	e.POST("/register", userH.PostRegister)
	e.POST("/register/username", userH.PostUsername)

	e.Logger.Fatal(e.Start(os.Getenv("ADDR")))
}
