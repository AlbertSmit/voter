package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	echo "github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
    return origin == "http://localhost:8080" || origin == "http://votevotevotevote.herokuapp.com"
	},
}

// WSHandler ...
func (a *App) WSHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	client := &Client{
		Hub:  a.hub,
		Conn: conn,
		Send: make(chan Message, 256),
	}

	a.hub.Register <- client

	client.Listen()

	return nil
}