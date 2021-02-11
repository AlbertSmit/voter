package main

//go:generate go run github.com/prisma/prisma-client-go generate

import (
	"context"
	"net/http"
	"strconv"

	db "github.com/albertsmit/voter/server/prisma-client"

	echo "github.com/labstack/echo/v4"
)

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

func createNewRoom(c echo.Context) error {
	client := db.NewClient()
	ctx := context.Background()

	if err := client.Prisma.Connect(); err != nil {
		return err
	}

	room, err := client.Room.CreateOne().Exec(ctx)

	if err != nil {
		return err
	}
	
	return c.JSONPretty(http.StatusOK, room, " ")
}
