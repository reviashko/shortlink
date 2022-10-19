package main

import (
	"os"
	"strconv"
	"sync"

	"github.com/reviashko/shortlink/pkg/app"
	"github.com/reviashko/shortlink/pkg/repository"

	"github.com/labstack/echo/v4"
)

var (
	acceptMutex = sync.Mutex{}
)

func main() {
	// Init echo server
	e := echo.New()

	DBURL := os.Getenv("DB_URL")
	db := repository.NewPostgreStorage(DBURL)

	AUTHDATA := os.Getenv("AUTH_DATA")
	if AUTHDATA == "" {
		AUTHDATA = `[{"Login":"joe", "Password": "secret"}]`
	}

	REFRESHTIME := os.Getenv("REFRESH_TIME")
	if REFRESHTIME == "" {
		REFRESHTIME = `10`
	}

	refreshTIME, err := strconv.Atoi(REFRESHTIME)
	if err != nil {
		panic(err)
	}

	web := app.NewWebServer(e, []byte(AUTHDATA))
	app.NewController(&db, &web, &acceptMutex, refreshTIME)

	// Start echo server
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	e.Logger.Fatal(e.Start(":" + PORT))

}
