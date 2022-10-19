package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/reviashko/shortlink/model"
	"github.com/reviashko/shortlink/pkg/app"
	"github.com/stretchr/testify/assert"
)

type mockStorage struct {
}

// Get func
func (s *mockStorage) Get() ([]model.URLItem, error) {

	return []model.URLItem{model.URLItem{ID: 1, Key: "key1", URL: "/url/111.jpg"}, model.URLItem{ID: 2, Key: "key2", URL: "/url/222.jpg"}}, nil
}

// GetSyncData func
func (s *mockStorage) GetSyncData(id int64) ([]model.URLItem, error) {
	return []model.URLItem{model.URLItem{ID: 3, Key: "key3", URL: "/url/333.jpg"}}, nil
}

// Delete func
func (s *mockStorage) Delete(id int64) error {
	return nil
}

// Save func
func (s *mockStorage) Save(model.URLItem) error {

	return nil
}

// Init func
func (s *mockStorage) Init() error {

	return nil
}

var testMutex sync.Mutex

func TestGetRedirectURL(t *testing.T) {

	e := echo.New()

	storage := mockStorage{}
	web := app.NewWebServer(e, []byte(`[{"Login":"joe", "Password": "secret"}]`))
	cntrl := app.NewController(&storage, &web, &testMutex, 5)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:key")
	c.SetParamNames("key")
	c.SetParamValues("key1")

	if assert.NoError(t, web.BaseHandler(c, cntrl.ShortURLHandler)) {
		assert.Equal(t, http.StatusMovedPermanently, rec.Code)
	}
}

func TestGetURLList(t *testing.T) {

	e := echo.New()

	storage := mockStorage{}
	web := app.NewWebServer(e, []byte(`[{"Login":"joe", "Password": "secret"}]`))
	cntrl := app.NewController(&storage, &web, &testMutex, 5)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/")

	if assert.NoError(t, web.BaseHandler(c, cntrl.ListURLsHandler)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestSyncData(t *testing.T) {

	e := echo.New()

	storage := mockStorage{}
	web := app.NewWebServer(e, []byte(`[{"Login":"joe", "Password": "secret"}]`))
	cntrl := app.NewController(&storage, &web, &testMutex, 5)

	testMutex.Lock()
	dataInit := len(cntrl.Data)
	testMutex.Unlock()
	time.Sleep(time.Duration(6) * time.Second)
	testMutex.Lock()
	dataFinish := len(cntrl.Data)
	testMutex.Unlock()

	if dataFinish == dataInit {
		t.Errorf("Wrong sync data. Waiting 1 item diff. Current: %d\n", dataFinish)
	}

}

func TestAddNewURL(t *testing.T) {

	e := echo.New()

	storage := mockStorage{}
	web := app.NewWebServer(e, []byte(`[{"Login":"joe", "Password": "secret"}]`))
	cntrl := app.NewController(&storage, &web, &testMutex, 5)

	f := make(url.Values)
	f.Set("url", "/superurl/test")

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, web.BaseHandler(c, cntrl.AddURLsHandler)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteURL(t *testing.T) {

	e := echo.New()

	storage := mockStorage{}
	web := app.NewWebServer(e, []byte(`[{"Login":"joe", "Password": "secret"}]`))
	cntrl := app.NewController(&storage, &web, &testMutex, 5)

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/api/:key")
	c.SetParamNames("key")
	c.SetParamValues("key1")

	lenBefore := len(cntrl.Data)

	if assert.NoError(t, web.BaseHandler(c, cntrl.DeleteURLHandler)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	lenAfter := len(cntrl.Data)

	if !((lenBefore - lenAfter) == 1) {
		t.Errorf("Wrong deleted items count\n")
	}
}
