// @see https://github.com/gorilla/websocket/issues/46#issuecomment-227906715
// and https://github.com/gorilla/websocket/blob/master/examples/chat/client.go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

const (
    // Time allowed to write a message to the peer.
    writeWait = 10 * time.Second

    // Time allowed to read the next pong message from the peer.
    pongWait = 60 * time.Second

    // Send pings to peer with this period. Must be less than pongWait.
    pingPeriod = (pongWait * 9) / 10

    // Maximum message size allowed from peer.
    maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
    // The websocket connection.
    ws *websocket.Conn

    // Buffered channel of outbound messages.
    send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump() {
    c := s.conn
    defer func() {
        h.unregister <- s
        c.ws.Close()
    }()
    c.ws.SetReadLimit(maxMessageSize)
    c.ws.SetReadDeadline(time.Now().Add(pongWait))
    c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
    for {
        _, msg, err := c.ws.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
                log.Printf("error: %v", err)
            }
            break
        }
        m := message{msg, s.room}
        h.broadcast <- m
    }
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
    c.ws.SetWriteDeadline(time.Now().Add(writeWait))
    return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {
    c := s.conn
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.ws.Close()
    }()
    for {
        select {
            case message, ok := <-c.send:
                if !ok {
                    c.write(websocket.CloseMessage, []byte{})
                    return
                }
                if err := c.write(websocket.TextMessage, message); err != nil {
                    return
                }

                c.write(websocket.TextMessage, message)

                // Add queued chat messages to the current websocket message.
                n := len(c.send)
                for i := 0; i < n; i++ {
                    c.write(websocket.TextMessage, newline)
                    c.write(websocket.TextMessage, <-c.send)
                }

                if err := c.ws.Close(); err != nil {
                    return
                }
            case <-ticker.C:
                if err := c.write(websocket.PingMessage, []byte{}); err != nil {
                    return
                }
        }
    }
}

// serveWs handles websocket requests from the peer.
func serveWs(c echo.Context) error {
    ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
    room := c.QueryParam("room")
    log.Println(room)
    if room == "" {
        return echo.NewHTTPError(http.StatusBadRequest, "Please provide a room to join.")
    }

    if err != nil {
        log.Println(err)
        return err
    }

    con := &connection{send: make(chan []byte, 256), ws: ws}
    sub := subscription{con, room}
    h.register <- sub
    go sub.writePump()
    sub.readPump()

    return nil
}