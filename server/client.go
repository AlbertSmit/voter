package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

// Client websocket client
type Client struct {
	Hub  		*Hub
	Conn 		*websocket.Conn
	Send 		chan Message
	Room		string
}

func (c *Client) Read() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		// log.Println("Client send Pong!")
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		message := Message{}
		err := c.Conn.ReadJSON(&message)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WS: Client error: %v", err)

				break
			}

			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.Printf("WS: Client error: %v", err)

				break
			}

			log.Printf("WS: Read message error: %v", err)
			continue
		}

		log.Printf("New incoming message: %s", message)
		c.Hub.Broadcast <- message
	}
}

func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Conn.WriteJSON(message)
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Client error send ping: %v", err)
				return
			}
		}
	}
}

// Listen ...
func (c *Client) Listen() {
	go c.Read()
	go c.Write()
}