package main

import (
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

var (
	rooms       	=	make(map[string]map[Subscription]*Client)
	votes					= make(map[string]map[Subscription]*Vote)
	state					= make(map[string]*StatefulRoom)

	register 			= make(chan Subscription)
	unregister 		= make(chan Subscription)
	
	broadcast 		= make(chan Message)
	status				= make(chan Status)
	update				= make(chan Update)
	vote					= make(chan CastVote)
) 

// Run it.
func (h* Hub) Run() {
	for {
		select {
			case connection := <-register:
				connections := rooms[connection.room]
				role := provideRole(connections)

				if connections == nil {
					connections = make(map[Subscription]*Client)
					rooms[connection.room] = connections

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
				e := createTypedResponse("update", clients)

				// Send new subs around.
				for c := range connections {
					writeToClient(c.connection, e)
				}

			case message := <-status:
				connections := rooms[message.Sub.room]

				e := createTypedResponse("status", State{message.State.Status})

				for c := range connections {
					if err := writeToClient(c.connection, e); err != nil {
						terminateClient(c.connection)
					}
				}

			case message := <-broadcast:
				connections := rooms[message.Sub.room]

				e := createTypedResponse("message", Payload{
					From: message.Data.From,
					Message: message.Data.Message,
				})

				for c := range connections {
					if err := writeToClient(c.connection, e); err != nil {
						terminateClient(c.connection)
					}
				}

			case v := <-vote:
				client := rooms[v.Sub.room][*v.Sub]
				ticket := &Vote{
					Motivation: v.Data.Motivation,
					For: v.Data.For,
					From: client,
				}

				vs := votes[v.Sub.room]
				if vs == nil {
					vs = make(map[Subscription]*Vote)
					votes[v.Sub.room] = vs
				} 

				votes[v.Sub.room][*v.Sub] = ticket

				vts := getVotes(votes, v.Sub.room)
				e := createTypedResponse("vote", vts)

				// Send vote around.
				connections := rooms[v.Sub.room]
				for c := range connections {
					writeToClient(c.connection, e)
				}

			case update := <-update:
				user := rooms[update.Sub.room][*update.Sub]
				rooms[update.Sub.room][*update.Sub] = &Client{
					Name: update.Data.Name,
					UUID: user.UUID,
				}

				clients := getClients(rooms, update.Sub.room)
				e := createTypedResponse("update", clients)

				// Send new subs around.
				connections := rooms[update.Sub.room]
				for c := range connections {
					writeToClient(c.connection, e)
				}

			case subscription := <-unregister:
				connections := rooms[subscription.room]
				if connections != nil {
					if _, ok := connections[subscription]; ok {
						delete(connections, subscription)

						clients := getClients(rooms, subscription.room)
						e := createTypedResponse("update", clients)

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