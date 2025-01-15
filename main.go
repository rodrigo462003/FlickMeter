package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	// Create a new Echo instance
	e := echo.New()

	// Serve static files (index.html, styles.css, etc.) from the 'public' directory
	e.Static("/", "./")

	// Start the Echo server
	e.Logger.Fatal(e.Start(":8080"))
}

