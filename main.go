package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/rodrigo462003/FlickMeter/handlers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Static("/public", "./public")
	e.GET("/", handlers.HomeHandler)
	e.GET("/signIn", handlers.SignInGetHandler)
	e.POST("/signIn", handlers.SignInPostHandler)
	e.Logger.Fatal(e.Start(os.Getenv("ADDR")))
}
