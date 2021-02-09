module voter

// +heroku goVersion go1.15
go 1.15

// +heroku install ./prisma-client
require (
	github.com/joho/godotenv v1.3.0
	github.com/labstack/echo/v4 v4.1.17
	github.com/prisma/prisma-client-go v0.4.0
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
)
