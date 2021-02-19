package main

import (
	"encoding/json"
	"fmt"
)

// Helper to add subscription to Update struct.
func withUpdate(update Update, sub *Subscription) Update {
	update.Sub = sub
	return update
}

// getClients maps over clients in a room.
func getClients(rooms map[string]map[Subscription]*Client, room string ) []*Client {
	var clients []*Client
	for _, client := range rooms[room] {
		clients = append(clients, client)
	}

	return clients
}

// writeTypedResponse reduces json Marshal boilerplate.
func writeTypedResponse(messageType string, data interface{}) []byte {
	payload := &ReponseWithType{
		Type: messageType,
		Data: data,
	}

	e, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return e
}

// provideRole
func provideRole(connections map[Subscription]*Client) Role {
	var role Role
	if connections == nil {
		role = Admin
	} else {
		role = User
	}
	
	return role
}