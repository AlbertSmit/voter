package main

import (
	"fmt"
	"context"
	"strconv"
	"net/http"
	"encoding/json"
	"voter/prisma/db"
	
	"golang.org/x/net/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func getRoot(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

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

func getPost(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, Post!")
}

func getSinglePost(c echo.Context) error {
	client := db.NewClient()
	ctx := context.Background()
	id, _ := strconv.Atoi(c.Param("id"))

	post, err := client.Post.FindUnique(
		db.Post.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
			return err
	}

	result, _ := json.MarshalIndent(post, "", "  ")
	return c.JSON(http.StatusOK, result)
}

func postNewPost(c echo.Context) error {
	client := db.NewClient()
	ctx := context.Background()

	createdPost, err := client.Post.CreateOne(
		db.Post.Title.Set("Hi from Prisma!"),
		db.Post.Published.Set(true),
		db.Post.Desc.Set("Prisma is a database toolkit and makes databases easy."),
	).Exec(ctx)

	if err != nil {
			return err
	}

	result, _ := json.MarshalIndent(createdPost, "", "  ")
	return c.JSON(http.StatusOK, result)
}

//go:generate go run github.com/prisma/prisma-client-go generate

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", getRoot)
	e.GET("/ws", getWS)
	e.GET("/post", getPost)
	e.GET("/post/:id", getSinglePost)
	e.POST("/post", postNewPost)

	e.Logger.Fatal(e.Start(":1323"))
}