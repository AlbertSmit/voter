// Convert to https://github.com/gofiber/recipes/tree/master/clean-architecture
// and this https://medium.com/@apzuk3/input-validation-in-golang-bc24cdec1835
package main

import (
	"encoding/json"
	"fmt"
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
	Fiber 				*fiber.App
}

// Payload is being sent by the client.
type Payload struct {
	Type 					string `json:"type" validate:"required"`
	From					string `json:"from" validate:"required"`
	Message 			string `json:"message" validate:"required"`
}

// Message gets send around.
type Message struct {
	Data 					Payload
	Room 					string
}

// Client uses the service.
type Client struct{
	UUID					string `json:"type" validate:"required,uuid4"`
} 

// Subscription exist when you connect.
type Subscription struct {
	connection 		*websocket.Conn
	room 					string
}

var (
	rooms       	=	make(map[string]map[*websocket.Conn]bool)
	clients 			= make(map[*Subscription]Client)
	register 			= make(chan Subscription)
	broadcast 		= make(chan Message)
	unregister 		= make(chan Subscription)
) 

func main() {
	app := App{}

	app.Initialize()
	app.Run()
}

// Initialize the Fiber server.
func (a *App) Initialize() {
	a.loadEnv()	

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
		s := Subscription{c, c.Params("room")}

		defer func() {
			unregister <- s
			c.Close()
		}()

		register <- s

		for {
			var result Message
			if err := s.connection.ReadJSON(&result); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("Read error:", err)
				}

				return // Calls the deferred function, i.e. closes the connection on error
			}

			log.Printf("From: %s, Message: %s, Type: %s", result.Data.From, result.Data.Message, result.Data.Type)
			broadcast <- Message{result.Data, c.Params("room")}
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
			// log.Println("Message received:", message)
			connections := rooms[message.Room]
			for c := range connections {
				// Stringify the data.
				emp := &Payload{
					Type: message.Data.Type,
					From: message.Data.From,
					Message: message.Data.Message,
				}
				e, err := json.Marshal(emp)
				if err != nil {
						fmt.Println(err)
						return
				}

				// Send to clients.
				if err := c.WriteMessage(websocket.TextMessage, []byte(e)); err != nil {
					log.Println("write error:", err)

					s := Subscription{c, c.Params("room")}
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

func (a *App) loadEnv() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}