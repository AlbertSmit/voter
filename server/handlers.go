package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

// createNewRoom to join.
func createNewRoom(ctx *fiber.Ctx) error {
	uuid := uuid.NewString()
	ctx.JSON(uuid)

	return nil
}

// When ANY to WS-endpoint.
func upgradeToWebSocket(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) { 
		return c.Next()
	}
	
	return c.SendStatus(fiber.StatusUpgradeRequired)
}

// When initializing WS.
func handleWebSocket(c *websocket.Conn) {
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

			return 
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
}