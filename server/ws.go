package main

import (
	"log"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	echo "github.com/labstack/echo/v4"
)

func serveSocket(c echo.Context) error {
	conn, _, _, err := ws.UpgradeHTTP(c.Request(), c.Response())
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer conn.Close()

	for {
		// Read
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			return err
		}

		// Write
		err = wsutil.WriteServerMessage(conn, op, msg)
		if err != nil {
			return err
		}
	}
}