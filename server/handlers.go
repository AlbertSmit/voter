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
	ctx.Response().Header.Add("X-Super-Admin", "Absolutely!")
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
				
				msg.Sub = &s
				broadcast <- msg
				
			case "status":
				var sts Status
				err := mapstructure.Decode(result, &sts)
				if err != nil {
					return
				}
				
				sts.Sub = &s
				status <-	sts

			case "update":
				var updt Update
				err := mapstructure.Decode(result, &updt)
				if err != nil {
					return
				}
				
				updt.Sub = &s
				update <- updt

			case "vote":
				var vt CastVote
				err := mapstructure.Decode(result, &vt)
				if err != nil {
					return
				}
				
				vt.Sub = &s
				vote <- vt

			case "control":
				var ctrl Control
				err := mapstructure.Decode(result, &ctrl)
				if err != nil {
					return
				}
				
				ctrl.Sub = &s
				control <- ctrl
			}
		}
}