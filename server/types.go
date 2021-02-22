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
	Sub						*Subscription
}

// State is the Status payload
type State struct {
	Status				string `json:"status"`
}

// Status for a room.
type Status struct {
	Type 					string `json:"type" validate:"required"`
	State 				State
	Sub						*Subscription
}

// Details to perform.
type Details struct {
	Name					string `json:"name"`
}

// Update gets send around.
type Update struct {
	Type 					string `json:"type" validate:"required"`
	Data 					Details
	Sub						*Subscription
}

// Client uses the service.
type Client struct {
	UUID					string `json:"uuid" validate:"required,uuid4"`
	Name					string `json:"name"`
	Role					Role `json:"role"`
} 

// Subscription exist when you connect.
type Subscription struct {
	connection 		*websocket.Conn
	room 					string
}

// Role for a user.
type Role int
const (
	// Admin rules supreme.
	Admin Role = iota
	// User follows.
	User
)

func (r Role) String() string {
	return [...]string{"Admin", "User"}[r]
}

// StatefulRoom holds room state.
type StatefulRoom struct {
	State					string `json:"state"`
	Pointer				int `json:"pointer"`
}

// Initial state of a room
type Initial struct {
	State					StatefulRoom 
	Clients				[]*Client
}

// Vote for voting.
type Vote struct {
	From					*Client
	For						*Client `json:"for"`
	Motivation		string `json:"motivation"`
}

// CastVote that is cast.
type CastVote struct {
	Type 					string `json:"type" validate:"required"`
	Data 					Vote
	Sub						*Subscription
}

// Command to give out
type Command struct {
	Pointer				int
}

// Control server
type Control struct {
	Type 					string `json:"type" validate:"required"`
	Data 					Command
	Sub						*Subscription
}