package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

// Run it.
func (h* Hub) Run() {
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