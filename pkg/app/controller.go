package app

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/reviashko/shortlink/model"
	"github.com/reviashko/shortlink/pkg/repository"
)

// Controller struct
type Controller struct {
	Storage    repository.StorageInterface
	Web        WebServerInterface
	Data       model.SafeMap
	SortedKeys model.SafeStringArray
	MaxID      int64
}

// NewController func
func NewController(storage repository.StorageInterface, web WebServerInterface, mutex *sync.Mutex, refreshTimeSec int) Controller {
	instance := Controller{Storage: storage, Web: web, Data: model.SafeMap{Mx: mutex}, SortedKeys: model.SafeStringArray{Mx: mutex}, MaxID: 0}
	instance.init(refreshTimeSec)
	return instance
}

// SyncData func
func (c *Controller) SyncData(refreshTimeSec int) {

	for {
		time.Sleep(time.Duration(refreshTimeSec) * time.Second)

		data, err := c.Storage.GetSyncData(c.MaxID)
		if err != nil {
			fmt.Printf("UpdateDataDiff err: %s", err.Error())
		}

		// TODO: embrace deleted items
		for _, item := range data {
			c.Data.Add(item.Key, model.ShortURLItem{ID: item.ID, URL: item.URL})
			c.SortedKeys.Append(item.Key, false)
		}
		c.SortedKeys.Sort()
	}
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

// init func
func (c *Controller) init(refreshTimeSec int) {

	err := c.initStorage()
	if err != nil {
		panic(err.Error())
	}

	c.initHandlers()
	go c.SyncData(refreshTimeSec)
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

	c.Data.Init()
	c.SortedKeys.Init(len(data))

	for _, item := range data {
		c.Data.Add(item.Key, model.ShortURLItem{ID: item.ID, URL: item.URL})
		c.SortedKeys.Append(item.Key, false)
		if item.ID > c.MaxID {
			c.MaxID = item.ID
		}
	}

	c.SortedKeys.Sort()

	return nil
}

// ShortURLHandler func
func (c *Controller) ShortURLHandler(key string, url string) (int, string, string) {

	item, exists := c.Data.Get(key)
	if !exists {
		return http.StatusNotFound, "Not found", model.StringResp
	}

	return http.StatusMovedPermanently, item.URL, model.RedirectResp
}

// ListURLsHandler func
func (c *Controller) ListURLsHandler(key string, url string) (int, string, string) {
	var sb strings.Builder

	rowTemplate, _ := os.ReadFile("templates/row.html")
	row := string(rowTemplate)

	for _, key := range c.SortedKeys.Get() {
		item, _ := c.Data.Get(key)
		sb.WriteString(fmt.Sprintf(row, key, key, item.URL, key))
	}

	file, _ := os.ReadFile("templates/index.html")
	content := string(file)
	content = strings.Replace(content, "{{%content%}}", sb.String(), -1)

	return http.StatusOK, content, model.HTMLResp
}

// AddURLsHandler func
func (c *Controller) AddURLsHandler(key string, url string) (int, string, string) {

	if len(key) == 0 {
		key = c.NewID()
		keyIsExists := true
		loops := 0
		for keyIsExists {

			_, keyIsExists = c.Data.Get(key)
			if keyIsExists {
				key = c.NewID()
			}

			if loops > 50 {
				return http.StatusInternalServerError, "no free keys", model.JSONResp
			}
			loops++
		}
	}

	if url == "" {
		return http.StatusInternalServerError, "error", model.JSONResp
	}

	if c.Data.IsExists(url) {
		return http.StatusOK, "ok", model.JSONResp
	}

	newID := time.Now().Unix()
	err := c.Storage.Save(model.URLItem{ID: newID, Key: key, URL: url})
	if err != nil {
		return http.StatusInternalServerError, "error", model.JSONResp
	}

	c.Data.Add(key, model.ShortURLItem{ID: newID, URL: url})
	c.SortedKeys.Append(key, true)

	return http.StatusOK, "ok", model.JSONResp
}

// LoginHandler func
func (c *Controller) LoginHandler(key string, url string) (int, string, string) {
	file, _ := os.ReadFile("templates/login.html")
	return http.StatusOK, string(file), model.HTMLResp
}

// DeleteURLHandler func
func (c *Controller) DeleteURLHandler(key string, url string) (int, string, string) {

	if key == "" {
		return http.StatusInternalServerError, "error", model.JSONResp
	}

	item, isExists := c.Data.Get(key)
	if !isExists {
		return http.StatusOK, "ok", model.JSONResp
	}

	err := c.Storage.Delete(item.ID)
	if err != nil {
		return http.StatusInternalServerError, "error", model.JSONResp
	}

	c.Data.Delete(key)
	c.SortedKeys.Delete(key)

	return http.StatusOK, "ok", model.JSONResp
}

// initHandlers func
func (c *Controller) initHandlers() {

	c.Web.AddGetHandler("/", c.ListURLsHandler)
	c.Web.AddUGetHandler("/:key", c.ShortURLHandler)
	c.Web.AddGetHandler("/login", c.LoginHandler)
	c.Web.AddPostHandler("/api", c.AddURLsHandler)
	c.Web.AddDeleteHandler("/api/:key", c.DeleteURLHandler)
}
