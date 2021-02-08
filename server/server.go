package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"voter/prisma/db"

	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	result, _ := json.MarshalIndent(post, "", "  ")
	return c.JSON(http.StatusOK, result)
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

	result, _ := json.MarshalIndent(createdPost, "", "  ")
	return c.JSON(http.StatusOK, result)
}

func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	e := echo.New()

	e.Use(middleware.Logger())

	e.POST("/post", postNewPost)
	e.GET("/post/:id", getSinglePost)

	e.Logger.Fatal(e.Start(":1323"))
}