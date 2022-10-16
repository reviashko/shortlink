package app

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/reviashko/shortlink/model"
	"github.com/reviashko/shortlink/pkg/repository"
)

// Controller struct
type Controller struct {
	Storage repository.StorageInterface
	Web     WebServerInterface
	Data    map[string]string
	Mutex   *sync.RWMutex
}

// NewController func
func NewController(storage repository.StorageInterface, web WebServerInterface, mutex *sync.RWMutex) Controller {
	instance := Controller{Mutex: mutex, Storage: storage, Web: web, Data: map[string]string{}}
	instance.init()
	return instance
}

// NewID func
func (c *Controller) NewID() string {

	var sb strings.Builder
	dataSet := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	loop := 0
	for loop < 6 {
		sb.WriteString(string(dataSet[rand.Intn(len(dataSet)-1)]))
		loop++
	}

	return sb.String()
}

// IsURLExist func
func (c *Controller) IsURLExist(url string) bool {

	for _, val := range c.Data {
		if val == url {
			return true
		}
	}

	return false
}

// init func
func (c *Controller) init() {
	err := c.initStorage()
	if err != nil {
		panic(err.Error())
	}

	c.initHandlers()
}

// initStorage func
func (c *Controller) initStorage() error {
	err := c.Storage.Init()
	if err != nil {
		return err
	}

	data, err := c.Storage.Get()
	if err != nil {
		return err
	}

	for _, item := range data {
		c.Data[item.Key] = item.URL
	}

	return nil
}

// initHandlers func
func (c *Controller) initHandlers() error {

	c.Web.AddGetHandler("/", func(key string, url string) (int, string, string) {

		var sb strings.Builder

		rowTemplate, _ := os.ReadFile("templates/row.html")
		row := string(rowTemplate)

		keys := make([]string, 0, len(c.Data))

		c.Mutex.Lock()
		for key := range c.Data {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			url := c.Data[key]
			sb.WriteString(fmt.Sprintf(row, key, key, url, key))
		}
		c.Mutex.Unlock()

		file, _ := os.ReadFile("templates/index.html")
		content := string(file)
		content = strings.Replace(content, "{{%content%}}", sb.String(), -1)

		return http.StatusOK, content, model.HTMLResp
	})

	c.Web.AddUGetHandler("/:key", func(key string, url string) (int, string, string) {

		c.Mutex.Lock()
		url, exists := c.Data[key]
		c.Mutex.Unlock()

		if !exists {
			return http.StatusNotFound, "Not found", model.StringResp
		}

		return http.StatusMovedPermanently, url, model.RedirectResp
	})

	c.Web.AddGetHandler("/login", func(key string, url string) (int, string, string) {
		file, _ := os.ReadFile("templates/login.html")
		return http.StatusOK, string(file), model.HTMLResp
	})

	c.Web.AddPostHandler("/api", func(key string, url string) (int, string, string) {

		if len(key) == 0 {
			key = c.NewID()
		}

		if url == "" {
			return http.StatusInternalServerError, "error", model.JSONResp
		}

		c.Mutex.Lock()
		isExists := c.IsURLExist(url)
		c.Mutex.Unlock()

		if isExists {
			return http.StatusOK, "ok", model.JSONResp
		}

		err := c.Storage.Save(model.URLItem{Key: key, URL: url})
		if err != nil {
			return http.StatusInternalServerError, "error", model.JSONResp
		}

		c.Data[key] = url

		return http.StatusOK, "ok", model.JSONResp
	})

	c.Web.AddDeleteHandler("/api/:key", func(key string, url string) (int, string, string) {

		if key == "" {
			return http.StatusInternalServerError, "error", model.JSONResp
		}

		c.Mutex.Lock()
		_, isExists := c.Data[key]
		c.Mutex.Unlock()
		if !isExists {
			return http.StatusOK, "ok", model.JSONResp
		}

		err := c.Storage.Delete(key)
		if err != nil {
			return http.StatusInternalServerError, "error", model.JSONResp
		}

		delete(c.Data, key)

		return http.StatusOK, "ok", model.JSONResp
	})

	return nil
}
