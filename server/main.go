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

type message struct {
	data 					string
	room 					string
}

type client struct{} 
type subscription struct {
	connection 		*websocket.Conn
	room 					string
}

var (
	rooms       	=	make(map[string]map[*websocket.Conn]bool)
	clients 			= make(map[*subscription]client)
	register 			= make(chan subscription)
	broadcast 		= make(chan message)
	unregister 		= make(chan subscription)
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
	go runHub()
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
	ws.Get("/:room", websocket.New(func(c *websocket.Conn) {
		s := subscription{c, c.Params("room")}

		defer func() {
			unregister <- s
			c.Close()
		}()

		register <- s

		for {
			messageType, msg, err := s.connection.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("Read error:", err)
				}

				return // Calls the deferred function, i.e. closes the connection on error
			}

			if messageType == websocket.TextMessage {
				broadcast <- message{string(msg), c.Params("room")}
				log.Println("Websocket message received of type text", messageType)
			} else {
				log.Println("Websocket message received of type", messageType)
			}
		}
	}))
}

func runHub() {
	for {
		select {
		case connection := <-register:
			connections := rooms[connection.room]
			if connections == nil {
					connections = make(map[*websocket.Conn]bool)
					rooms[connection.room] = connections
			}
			rooms[connection.room][connection.connection] = true

			log.Println("Connection registered")

		case message := <-broadcast:
			log.Println("Message received:", message)

			connections := rooms[message.room]
			for c := range connections {
				if err := c.WriteMessage(websocket.TextMessage, []byte(message.data)); err != nil {
					log.Println("write error:", err)

					s := subscription{c, c.Params("room")}
					unregister <- s
					c.WriteMessage(websocket.CloseMessage, []byte{})
					c.Close()
				}
			}

		case subscription := <-unregister:
			connections := rooms[subscription.room]
			if connections != nil {
					if _, ok := connections[subscription.connection]; ok {
							delete(connections, subscription.connection)
							// subscription.connection.Close()
							if len(connections) == 0 {
									delete(rooms, subscription.room)
							}
					}
			}

			log.Println("Connection unregistered")
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