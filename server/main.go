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
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
)

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
	hub := Hub{}

	app.Initialize()
	go hub.Run()
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
}

func createNewRoom(ctx *fiber.Ctx) error {
	uuid := uuid.NewString()
	ctx.JSON(uuid)

	return nil
}

// Helper to add subscription to Update struct.
func withUpdate(update Update, sub *Subscription) Update {
	update.Sub = sub
	return update
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

// Run the app
func (a *App) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	a.Fiber.Listen(":" + port)
}

// Load the env.
func (a *App) loadEnv() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}