package main

import (
	"github.com/gorilla/websocket"

	echo "github.com/labstack/echo/v4"
)

// func serveSocket(c echo.Context) error {
// 	conn, _, _, err := ws.UpgradeHTTP(c.Request(), c.Response())
// 	if err != nil {
// 		log.Fatal(err)
// 		return err
// 	}

// 	defer conn.Close()

// 	for {
// 		// Read
// 		msg, op, err := wsutil.ReadClientData(conn)
// 		if err != nil {
// 			return err
// 		}

// 		// Write
// 		err = wsutil.WriteServerMessage(conn, op, msg)
// 		if err != nil {
// 			return err
// 		}
// 	}
// }

var upgrader = websocket.Upgrader{}

// WSHandler ...
func (a *App) WSHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
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