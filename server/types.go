package main

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// App houses Fiber.
type App struct {
	Fiber 				*fiber.App
}

// Hub houses connections.
type Hub struct {}

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