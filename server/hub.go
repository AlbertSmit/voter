package main

import (
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

var (
	rooms       	=	make(map[string]map[Subscription]*Client)
	state					= make(map[string]*StatefulRoom)

	register 			= make(chan Subscription)
	unregister 		= make(chan Subscription)
	
	broadcast 		= make(chan Message)
	status				= make(chan Status)
	update				= make(chan Update)
) 

// Run it.
func (h* Hub) Run() {
	for {
		select {
			case connection := <-register:
				// Get room
				connections := rooms[connection.room]
				role := provideRole(connections)

				if connections == nil {
					// Create room
					connections = make(map[Subscription]*Client)
					rooms[connection.room] = connections

					// Set initial state
					state[connection.room] = &StatefulRoom{
						State: "WAITING",
					}
				} 

				// Add client to room
				rooms[connection.room][connection] = &Client{
					UUID: uuid.NewString(), 					
					Role: Role(role),
				}

				clients := getClients(rooms, connection.room)
				e := writeTypedResponse("update", clients)

				// Send new subs around.
				for c := range connections {
					c.connection.WriteMessage(websocket.TextMessage, []byte(e))
					c.connection.Close()
				}

			case message := <-status:
				connections := rooms[message.Room]

				e := writeTypedResponse("status", State{message.State.Status})

				for c := range connections {
					if err := c.connection.WriteMessage(websocket.TextMessage, []byte(e)); err != nil {
						s := Subscription{c.connection, c.connection.Params("room")}
						unregister <- s
						
						c.connection.WriteMessage(websocket.CloseMessage, []byte{})
						c.connection.Close()
					}
				}

			case message := <-broadcast:
				connections := rooms[message.Room]

				e := writeTypedResponse("message", Payload{
					From: message.Data.From,
					Message: message.Data.Message,
				})

				for c := range connections {
					if err := c.connection.WriteMessage(websocket.TextMessage, []byte(e)); err != nil {
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

				clients := getClients(rooms, update.Room)
				e := writeTypedResponse("update", clients)

				// Send new subs around.
				connections := rooms[update.Room]
				for c := range connections {
					c.connection.WriteMessage(websocket.TextMessage, []byte(e))
					c.connection.Close()
				}

			case subscription := <-unregister:
				connections := rooms[subscription.room]
				if connections != nil {
					if _, ok := connections[subscription]; ok {
						delete(connections, subscription)

						clients := getClients(rooms, subscription.room)
						e := writeTypedResponse("update", clients)

						// Notify other users of abscense.
						for c := range connections {
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