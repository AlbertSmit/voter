// @see https://github.com/tegalan/echo-ws/
package main

import "log"

// Hub ...
type Hub struct {
	// Clients 			map[*Client]bool
	Rooms 				map[string]map[*Client]bool

	Broadcast 		chan Message

	Register   		chan *Client
	Unregister 		chan *Client
}

// NewHub ...
func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		// Clients:    make(map[*Client]bool),
		Rooms:			make(map[string]map[*Client]bool),
	}
}

// Run the hub
func (hub *Hub) Run() {
	log.Println("WS Hub running.")

	for {
		select {
		case client := <-hub.Register:
			// New
			connections := hub.Rooms[client.Room]
			if connections == nil {
					connections = make(map[*Client]bool)
					hub.Rooms[client.Room] = connections
			}
			hub.Rooms[client.Room][client] = true

			// Old
			// h.Clients[client] = true
			log.Println("Client connected!")
		case client := <-hub.Unregister:
			// New
			connections := hub.Rooms[client.Room]
			if connections != nil {
				if _, ok := connections[client]; ok {
					delete(connections, client)
					close(client.Send)
					if len(connections) == 0 {
							delete(hub.Rooms, client.Room)
					}
				}
			}

			log.Println("Client disconnected!")

			// Old
			// if _, ok := h.Clients[client]; ok {
			// 	delete(h.Clients, client)
			// 	close(client.Send)

			// 	log.Println("Client disconnected!")
			// }
		case message := <-hub.Broadcast:
			// New
			connections := hub.Rooms[message.Room]
			for connection := range connections {
				select {
					case connection.Send <- message:
						log.Printf("Broadcast message: %s", message)
					default:
						close(connection.Send)
						delete(connections, connection)
						if len(connections) == 0 {
								delete(hub.Rooms, message.Room)
						}
				}
			}

			// Old
			// for client := range h.Clients {
			// 	select {
			// 	case client.Send <- message:
			// 		log.Printf("Broadcast message: %s", message)
			// 	default:
			// 		close(client.Send)
			// 		delete(h.Clients, client)
			// 	}
			// }
		}
	}
}