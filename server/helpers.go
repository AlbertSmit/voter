package main

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/websocket/v2"
)

// getClients maps over clients in a room.
func listRooms(rooms map[string]map[Subscription]*Client) []string {
	var rms []string
	for key := range rooms {
		rms = append(rms, key)
	}

	return rms
}

// getClients maps over clients in a room.
func getClients(rooms map[string]map[Subscription]*Client, room string) []*Client {
	var clients []*Client
	for _, client := range rooms[room] {
		clients = append(clients, client)
	}

	return clients
}

// getVotes maps over votes.
func getVotes(votes map[string]map[Subscription]*Vote, room string) []*Vote {
	var vts []*Vote
	for _, vote := range votes[room] {
		vts = append(vts, vote)
	}

	return vts
}

// createTypedResponse reduces json Marshal boilerplate.
func createTypedResponse(messageType string, data interface{}) []byte {
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

// write to clients
func writeWithoutTermination(c *websocket.Conn, message []byte) error {
	writeErr := c.WriteMessage(websocket.TextMessage, []byte(message))
	if (writeErr != nil) {
		return writeErr
	}

	return nil
}

// write to clients
func writeToClient(c *websocket.Conn, message []byte) error {
	writeErr := c.WriteMessage(websocket.TextMessage, []byte(message))
	if (writeErr != nil) {
		return writeErr
	}

	closeErr := c.Close()
	if (closeErr != nil) {
		return closeErr
	}

	return nil
}

// terminate connection
func terminateClient(c *websocket.Conn) {
	s := Subscription{c, c.Params("room")}
	unregister <- s
	
	c.WriteMessage(websocket.CloseMessage, []byte{})
	c.Close()
}