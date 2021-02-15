// Convert to https://github.com/gofiber/recipes/tree/master/clean-architecture
package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// App houses Fiber.
type App struct {
	Fiber *fiber.App
}

type client struct{} 

var (
	clients 			= make(map[*websocket.Conn]client)
	register 			= make(chan *websocket.Conn)
	broadcast 		= make(chan string)
	unregister 		= make(chan *websocket.Conn)
) 


func main() {
	app := App{}

	app.Initialize()
	app.Run()
}

// Initialize the Fiber server.
func (a *App) Initialize() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	app := fiber.New()
	a.Fiber = app

	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(compress.New())

	if os.Getenv("APP_ENV") == "production" {
		app.Get("/", func(c *fiber.Ctx) error {
			return c.Redirect("/web")
		})
		app.Static("/web", "./web")

		app.Get("/web/*", func(ctx *fiber.Ctx) error {
			return ctx.SendFile("./dist/index.html")
		})
	}

	a.InitRouter()
}

func createNewRoom(ctx *fiber.Ctx) error {
	uuid := uuid.NewString()
	ctx.JSON(uuid)

	return nil
}

// InitRouter all the routes.
func (a *App) InitRouter() {
	party := a.Fiber.Group("/api")
	party.Get("/room", createNewRoom)

	// WS Upgrade Middleware
	ws := party.Group("/socket", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { 
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	// WS
	ws.Get("/", websocket.New(func(c *websocket.Conn) {
		defer func() {
			unregister <- c
			c.Close()
		}()

		register <- c

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}

				return // Calls the deferred function, i.e. closes the connection on error
			}

			if messageType == websocket.TextMessage {
				broadcast <- string(message)
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}
	}))
}

func runHub() {
	for {
		select {
		case connection := <-register:
			clients[connection] = client{}
			log.Println("connection registered")

		case message := <-broadcast:
			log.Println("message received:", message)

			for connection := range clients {
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Println("write error:", err)

					unregister <- connection
					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
				}
			}

		case connection := <-unregister:
			delete(clients, connection)

			log.Println("connection unregistered")
		}
	}
}

// Run the app
func (a *App) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	a.Fiber.Listen(":" + port)
}