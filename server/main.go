package main

//go:generate go run github.com/prisma/prisma-client-go generate

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// App houses echo framework and WS hub.
type App struct {
	Echo *echo.Echo
	hub  *Hub
}

func main() {
	app := App{}

	app.Initialize()
	app.Run()
}

// Initialize the Echo server.
func (a *App) Initialize() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	e := echo.New()
	a.Echo = e

	a.hub = NewHub()

	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	if os.Getenv("APP_ENV") == "production" {
		e.Static("/", "./web")
	}

	a.InitRouter()
}

// InitRouter all the routes.
func (a *App) InitRouter() {
	e := a.Echo

	e.Any("/socket", serveSocket)
	e.POST("/post", postNewPost)
	e.GET("/post/:id", getSinglePost)
	e.POST("/room", createNewRoom)
}

// Run the app
func (a *App) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	go a.hub.Run()
	a.Echo.Logger.Fatal(a.Echo.Start(":" + port))
}