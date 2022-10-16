package main

import (
	"os"
	"sync"

	"github.com/reviashko/shortlink/pkg/app"
	"github.com/reviashko/shortlink/pkg/repository"

	"github.com/labstack/echo/v4"
)

var acceptMutex sync.RWMutex

func main() {
	// Init echo server
	e := echo.New()

	DBURL := os.Getenv("DB_URL")
	db := repository.NewPostgreStorage(DBURL)

	AuthData := os.Getenv("AUTH_DATA")
	if AuthData == "" {
		AuthData = `[{"Login":"joe", "Password": "secret"}]`
	}

	web := app.NewWebServer(e, []byte(AuthData))
	app.NewController(&db, &web, &acceptMutex)

	// Start echo server
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	e.Logger.Fatal(e.Start(":" + PORT))

}
