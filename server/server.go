package main

//go:generate go run github.com/prisma/prisma-client-go generate

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"voter/prisma/db"

	"github.com/joho/godotenv"
	"golang.org/x/net/websocket"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func getWS(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			// Write
			err := websocket.Message.Send(ws, "Hello, Client!")
			if err != nil {
				c.Logger().Error(err)
			}

			// Read
			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
			}
			fmt.Printf("%s\n", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func getSinglePost(c echo.Context) error {
	client := db.NewClient()
	ctx := context.Background()
	id, _ := strconv.Atoi(c.Param("id"))

	if err := client.Prisma.Connect(); err != nil {
		return err
	}

	post, err := client.Post.FindUnique(
		db.Post.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
			return err
	}

	return c.JSONPretty(http.StatusOK, post, "")
}

func postNewPost(c echo.Context) error {
	client := db.NewClient()
	ctx := context.Background()

	if err := client.Prisma.Connect(); err != nil {
		return err
	}

	createdPost, err := client.Post.CreateOne(
		db.Post.Title.Set("Hi from Prisma!"),
		db.Post.Published.Set(true),
		db.Post.Desc.Set("Prisma is a database toolkit and makes databases easy."),
	).Exec(ctx)

	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, createdPost, " ")
}

func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/ws", getWS)
	e.POST("/post", postNewPost)
	e.GET("/post/:id", getSinglePost)

	e.Logger.Fatal(e.Start(":1323"))
}