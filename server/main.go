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
	"github.com/mitchellh/mapstructure"
)

// App houses Fiber.
type App struct {
	Fiber 				*fiber.App
}

// IncomingRequest to catch
type IncomingRequest struct {
  Type 					string `json:"type"`
  Data 					json.RawMessage
}

// ReponseWithType to return to client.
type ReponseWithType struct {
	Type					string `json:"type"`
	Data					interface{} `json:"data"`
}

// Payload is being sent by the client.
type Payload struct {
	From					string `json:"from" validate:"required"`
	Message 			string `json:"message" validate:"required"`
}

// Message gets send around.
type Message struct {
	Type 					string `json:"type" validate:"required"`
	Data 					Payload
	Room 					string `json:"room"`
}

// State is the Status payload
type State struct {
	Status				string `json:"status"`
}

// Status for a room.
type Status struct {
	Type 					string `json:"type" validate:"required"`
	State 				State
	Room 					string `json:"room"`
}

// Details to perform.
type Details struct {
	Name					string `json:"name"`
}

// Update gets send around.
type Update struct {
	Type 					string `json:"type" validate:"required"`
	Data 					Details
	Room 					string `json:"room"`
	Sub						*Subscription
}

// Client uses the service.
type Client struct{
	UUID					string `json:"uuid" validate:"required,uuid4"`
	Name					string `json:"name"`
} 

// Subscription exist when you connect.
type Subscription struct {
	connection 		*websocket.Conn
	room 					string
}

var (
	// Manage clients
	rooms       	=	make(map[string]map[Subscription]*Client)

	// Manage connections
	register 			= make(chan Subscription)
	unregister 		= make(chan Subscription)
	
	// Send data
	broadcast 		= make(chan Message)
	status				= make(chan Status)
	update				= make(chan Update)
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

	// First, init the routes
	a.InitRouter()

	// Then init static 
	// (since we are on root)
	if os.Getenv("APP_ENV") == "production" {
		app.Static("/", "web")
		app.Get("/*", func(ctx *fiber.Ctx) error {
			return ctx.SendFile("./web/index.html")
		})
	}

	go runHub()
}

func createNewRoom(ctx *fiber.Ctx) error {
	uuid := uuid.NewString()
	ctx.JSON(uuid)

	return nil
}

// Helper
func withUpdate(updt Update, sub *Subscription) Update {
	updt.Sub = sub
	return updt
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
			// Interface for type switching.
			// Parse JSON from incoming message.
			var result map[string]interface{}
			if err := s.connection.ReadJSON(&result); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("Read error:", err)
				}

				return // Calls the deferred function, i.e. closes the connection on error
			}

			switch result["type"] {
				case "message":
					var msg Message
					err := mapstructure.Decode(result, &msg)
					if err != nil {
						return
					}
					
					broadcast <- msg
				case "status":
					var sts Status
					err := mapstructure.Decode(result, &sts)
					if err != nil {
						return
					}
					
					status <-	sts

				case "update":
					var updt Update
					err := mapstructure.Decode(result, &updt)
					if err != nil {
						return
					}
					
					update <- withUpdate(updt, &s)
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
					connections = make(map[Subscription]*Client)
					rooms[connection.room] = connections
				}
				rooms[connection.room][connection] = &Client{ UUID: uuid.NewString() }

				// Send new subs around.
				for c := range connections {
					clients := []*Client{}
					for _, client := range rooms[connection.room] {
						clients = append(clients, client)
					}

					// Emit the full room (?).
					payload := &ReponseWithType{
						Type: "update",
						Data: clients,
					}

					e, err := json.Marshal(payload)
					if err != nil {
							fmt.Println(err)
							return
					}

					c.connection.WriteMessage(websocket.TextMessage, []byte(e))
					c.connection.Close()
				}
				
			case message := <-status:
				connections := rooms[message.Room]
				for c := range connections {
					// Stringify the data.
					payload := &ReponseWithType{
						Type: "status",
						Data: State{message.State.Status},
					}
					e, err := json.Marshal(payload)
					if err != nil {
						fmt.Println(err)
						return
					}

					// Send to clients.
					if err := c.connection.WriteMessage(websocket.TextMessage, []byte(e)); err != nil {
						log.Println("write error:", err)

						s := Subscription{c.connection, c.connection.Params("room")}
						unregister <- s
						
						c.connection.WriteMessage(websocket.CloseMessage, []byte{})
						c.connection.Close()
					}
				}


			case message := <-broadcast:
				connections := rooms[message.Room]
				for c := range connections {
					// Stringify the data.
					payload := &ReponseWithType{
						Type: "message",
						Data: Payload{
							From: message.Data.From,
							Message: message.Data.Message,
						},
					}
					e, err := json.Marshal(payload)
					if err != nil {
							fmt.Println(err)
							return
					}

					// Send to clients.
					if err := c.connection.WriteMessage(websocket.TextMessage, []byte(e)); err != nil {
						log.Println("write error:", err)

						s := Subscription{c.connection, c.connection.Params("room")}
						unregister <- s

						c.connection.WriteMessage(websocket.CloseMessage, []byte{})
						c.connection.Close()
					}
				}

			case update := <-update:
				user := rooms[update.Room][*update.Sub]
				rooms[update.Room][*update.Sub] = &Client{
					Name: update.Data.Name,
					UUID: user.UUID,
				}

				// Send new subs around.
				connections := rooms[update.Room]
				for c := range connections {
					clients := []*Client{}
					for _, client := range rooms[update.Room] {
						clients = append(clients, client)
					}

					// Emit the full room (?).
					payload := &ReponseWithType{
						Type: "update",
						Data: clients,
					}

					e, err := json.Marshal(payload)
					if err != nil {
							fmt.Println(err)
							return
					}

					c.connection.WriteMessage(websocket.TextMessage, []byte(e))
					c.connection.Close()
				}

			case subscription := <-unregister:
				connections := rooms[subscription.room]
				if connections != nil {
					if _, ok := connections[subscription]; ok {
						delete(connections, subscription)

						// Notify other users of abscense.
						for c := range connections {
							clients := []*Client{}
							for _, client := range rooms[subscription.room] {
								log.Println("Client", client)
								clients = append(clients, client)
							}

							// Emit the full room (?).
							payload := &ReponseWithType{
								Type: "update",
								Data: clients,
							}

							e, err := json.Marshal(payload)
							if err != nil {
									fmt.Println(err)
									return
							}

							c.connection.WriteMessage(websocket.TextMessage, []byte(e))
							c.connection.Close()
						}

						if len(connections) == 0 {
							delete(rooms, subscription.room)
						}
					}
				}
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