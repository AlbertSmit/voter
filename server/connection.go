// @see https://github.com/gorilla/websocket/issues/46#issuecomment-227906715
// and https://github.com/gorilla/websocket/blob/master/examples/chat/client.go
// and https://github.com/umirode/echo-socket.io/blob/master/wrapper.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/labstack/echo/v4"
)

func serveSocket(c echo.Context) error {
    // server, err := socketio.NewServer(nil)
    allowOrigin := func(r *http.Request) bool {
        return true
    }

    server, err := socketio.NewServer(&engineio.Options{
        Transports: []transport.Transport{
            &polling.Transport{
                Client: &http.Client{
                    Timeout: time.Minute,
                },
                CheckOrigin: allowOrigin,
            },
            &websocket.Transport{
                CheckOrigin: allowOrigin,
            },
        },
    })

    if (err != nil) {
        log.Println(err)
        return err
    }

    server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())

		return nil
	})

    server.OnEvent("/", "room", func(s socketio.Conn, msg string) {
        s.Join(msg)
    })

    server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()
	defer server.Close()

    server.ServeHTTP(c.Response(), c.Request())
    return nil
}
