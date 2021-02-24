package main

import (
	"fmt"

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
	control				= make(chan Control)
	status				= make(chan Status)
	update				= make(chan Update)
	vote					= make(chan CastVote)
) 

// Run it.
func (h* Hub) Run() {
	for {
		select {
			/* 
				On register, run this.
			*/
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
					Role: role,
				}

				// Send out state
				a := createTypedResponse("status", state[connection.room])
				writeWithoutTermination(connection.connection, a)

				clients := getClients(rooms, connection.room)
				e := createTypedResponse("update", clients)

				// Send new subs around.
				for c := range connections {
					writeToClient(c.connection, e)
				}

				
			/* 
				On status, run this.
			*/
			case message := <-status:
				connections := rooms[message.Sub.room]
				
				state[message.Sub.room] = &StatefulRoom{
					State: message.State.Status,
				}

				e := createTypedResponse("status", State{message.State.Status})

				for c := range connections {
					if err := writeToClient(c.connection, e); err != nil {
						terminateClient(c.connection)
					}
				}

				
			/* 
				On message, run this.
			*/
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


			/* 
				On vote, run this.
			*/
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


			/* 
				On update, run this.
			*/
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


			/* 
				On control, run this.
			*/
			case control := <-control:
				/* 
					Only admins allowed here. 
					You shall not pass.
					GTFO y'all.
					Gotcha,
				*/
				if rooms[control.Sub.room][*control.Sub].Role != Admin {
					control.Sub.connection.Close()
				}

				pointer := control.Data.Pointer
				state[control.Sub.room].Pointer = pointer
				fmt.Println("Next pointer passed by admin:", pointer)

				connections := rooms[control.Sub.room]
				e := createTypedResponse("control", StatefulRoom{
					State: state[control.Sub.room].State,
					Pointer: pointer,
				})

				// Send new pointer around.
				for c := range connections {
					writeToClient(c.connection, e)
				}


			/* 
				On unregister, run this.
			*/
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