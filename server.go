package main

import (
	"context"
	"net/http"
	"encoding/json"
	"voter/prisma/db"
	"github.com/labstack/echo/v4"
)

func main() {
	// Prisma 2
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
			return
	}

	defer func() {
			if err := client.Prisma.Disconnect(); err != nil {
					panic(err)
			}
	}()

	ctx := context.Background()

	// Echo
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})


	e.GET("/post/:id", func(c echo.Context) error {
		id := c.Param("id")

		post, err := client.Post.FindUnique(
			db.Post.ID.Equals(id),
		).Exec(ctx)

		if err != nil {
				return err
		}

		result, _ := json.MarshalIndent(post, "", "  ")
		return c.JSON(http.StatusOK, result)
	})

	e.POST("/post", func(c echo.Context) error {
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
	})

	e.Logger.Fatal(e.Start(":1323"))
}