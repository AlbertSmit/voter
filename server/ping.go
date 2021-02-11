package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Ping ...
func (a *App) Ping(c echo.Context) error {
	a.hub.Broadcast <- Message{
		Type: "ping",
		Body: "Ping from server!",
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Pong!"})
}