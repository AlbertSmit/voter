package main

//go:generate go run github.com/prisma/prisma-client-go generate

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	db "github.com/albertsmit/voter/server/prisma-client"

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

func main() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://votevotevotevote.herokuapp.com"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	if os.Getenv("APP_ENV") == "production" {
		e.Static("/", "./web")
	}

	e.Any("/socket.io", serveSocket)

	e.POST("/post", postNewPost)
	e.GET("/post/:id", getSinglePost)
	e.POST("/room", createNewRoom)

	port := os.Getenv("PORT")
	if port == "" {
			log.Fatal("$PORT must be set")
	}

	e.Logger.Fatal(e.Start(":" + port))
}