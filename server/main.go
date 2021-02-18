// Convert to https://github.com/gofiber/recipes/tree/master/clean-architecture
// and this https://medium.com/@apzuk3/input-validation-in-golang-bc24cdec1835
package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
)

func main() {
	app := App{}
	hub := Hub{}

	app.Initialize()
	go hub.Run()

	app.Run()
}

// Initialize the Fiber server.
func (a *App) Initialize() {
	a.Fiber = fiber.New()
	
	a.EnvVars()	
	a.InitMiddlewares()
	a.InitRouter()
	a.InitSPA()
}

// InitRouter all the routes.
func (a *App) InitRouter() {
	// '/api' route
	party := a.Fiber.Group("/api")
	party.Get("/room", createNewRoom)

	// '/socket' route
	ws := party.Group("/socket", upgradeToWebSocket)
	ws.Get("/:room", 	websocket.New(func(c *websocket.Conn) {
		handleWebSocket(c)
	}))
}

// Run the app
func (a *App) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	a.Fiber.Listen(":" + port)
}

// EnvVars to load.
func (a *App) EnvVars() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

// InitMiddlewares for Fiber.
func (a *App) InitMiddlewares() {
	a.Fiber.Use(recover.New())
	a.Fiber.Use(compress.New())
	a.Fiber.Use(cors.New(cors.Config{
		ExposeHeaders: "X-Super-Admin",
		AllowOrigins: "http://localhost:8080, http://votevotevotevote.herokuapp.com",
	}))
}

// InitSPA to serve to client.
func (a *App) InitSPA() {
	if os.Getenv("APP_ENV") == "production" {
		a.Fiber.Static("/", "web")
		a.Fiber.Get("/*", func(ctx *fiber.Ctx) error {
			return ctx.SendFile("./web/index.html")
		})
	}
}