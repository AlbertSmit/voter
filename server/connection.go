// @see https://github.com/gorilla/websocket/issues/46#issuecomment-227906715
// and https://github.com/gorilla/websocket/blob/master/examples/chat/client.go
// and https://github.com/umirode/echo-socket.io/blob/master/wrapper.go
package main

import (
	"fmt"
	"log"

	socketio "github.com/googollee/go-socket.io"
	"github.com/labstack/echo/v4"
)

func serveSocket(c echo.Context) error {
    server, err := socketio.NewServer(nil)
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
