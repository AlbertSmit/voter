package main

import "log"

// Hub ...
type Hub struct {
	Clients map[*Client]bool
	Broadcast chan Message
	Register   chan *Client
	Unregister chan *Client
}

// NewHub ...
func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

// Run the hub
func (h *Hub) Run() {
	log.Println("WS Hub running.")

	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Println("Client connected!")
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				// Delete from array
				delete(h.Clients, client)
				// Close channel
				close(client.Send)

				log.Println("Client disconnected!")
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
					log.Printf("Broadcast message: %s", message)
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}

	}
}