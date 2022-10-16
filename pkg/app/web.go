package app

import (
	"encoding/json"

	"github.com/reviashko/shortlink/model"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// WebServerInterface interface
type WebServerInterface interface {
	AddGetHandler(route string, nextFn func(string, string) (int, string, string))
	AddUGetHandler(route string, nextFn func(string, string) (int, string, string))
	AddPostHandler(route string, nextFn func(string, string) (int, string, string))
	AddDeleteHandler(route string, nextFn func(string, string) (int, string, string))
}

// WebServer struct
type WebServer struct {
	Echo   *echo.Group
	UnSafe *echo.Echo
}

// NewWebServer func
func NewWebServer(echo *echo.Echo, authJSON []byte) WebServer {

	if len(authJSON) < 10 {
		panic("wrong authJSON")
	}

	var auth []model.AuthItem
	err := json.Unmarshal(authJSON, &auth)
	if err != nil {
		panic("Unmarshal authJSON error")
	}

	instance := WebServer{Echo: echo.Group(""), UnSafe: echo}
	instance.initAuth(auth)
	return instance
}

// InitAuth func
func (w *WebServer) initAuth(auth []model.AuthItem) {
	w.Echo.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {

		for _, item := range auth {
			if username == item.Login && password == item.Password {
				return true, nil
			}
		}
		return false, nil
	}))
}

func (w *WebServer) httpHandler(c echo.Context, route string, nextFn func(string, string) (int, string, string)) error {

	url := c.FormValue("url")
	key := c.FormValue("key")
	if key == "" {
		key = c.Param("key")
	}

	code, data, responsetype := nextFn(key, url)

	if responsetype == "Redirect" {
		return c.Redirect(code, data)
	}

	if responsetype == "String" {
		return c.String(code, data)
	}

	if responsetype == "JSON" {
		c.Response().Header().Set("Content-Type", "application/json")
		return c.JSON(code, map[string]string{"status": data})
	}

	return c.HTML(code, data)
}

// AddGetHandler func
func (w *WebServer) AddGetHandler(route string, nextFn func(string, string) (int, string, string)) {
	w.Echo.GET(route, func(c echo.Context) error {
		return w.httpHandler(c, route, nextFn)
	})
}

// AddUGetHandler func
func (w *WebServer) AddUGetHandler(route string, nextFn func(string, string) (int, string, string)) {
	w.UnSafe.GET(route, func(c echo.Context) error {
		return w.httpHandler(c, route, nextFn)
	})
}

// AddPostHandler func
func (w *WebServer) AddPostHandler(route string, nextFn func(string, string) (int, string, string)) {
	w.Echo.POST(route, func(c echo.Context) error {
		return w.httpHandler(c, route, nextFn)
	})
}

// AddDeleteHandler func
func (w *WebServer) AddDeleteHandler(route string, nextFn func(string, string) (int, string, string)) {
	w.Echo.DELETE(route, func(c echo.Context) error {
		return w.httpHandler(c, route, nextFn)
	})
}
