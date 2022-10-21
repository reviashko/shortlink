package repository

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/reviashko/shortlink/model"
)

// StorageInterface interface
type StorageInterface interface {
	Get() ([]model.URLItem, error)
	GetSyncData(id int64) ([]model.URLItem, error)
	Delete(id int64) error
	Save(model.URLItem) error
	Init() error
}

// PostgreStorage struct
type PostgreStorage struct {
	DB *sql.DB
}

// NewPostgreStorage func
func NewPostgreStorage(dbURL string) PostgreStorage {

	if dbURL == "" {
		panic("Wrong database url!")
	}
	dbURL, _ = pq.ParseURL(dbURL)
	db, _ := sql.Open("postgres", dbURL)

	fmt.Println("Try database connection...")

	err := db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("database connected")

	return PostgreStorage{DB: db}
}

// Init func
func (p *PostgreStorage) Init() error {

	_, err := p.DB.Exec("create table if not exists shortlinks(id bigint, key varchar, url varchar, constraint pk_shortlinks primary key(id))")
	if err != nil {
		return err
	}

	return nil
}

// Get func
func (p *PostgreStorage) Get() ([]model.URLItem, error) {

	result := make([]model.URLItem, 0, 50)

	rows, err := p.DB.Query("select id, key, url from shortlinks")
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		item := model.URLItem{}
		rows.Scan(&item.ID, &item.Key, &item.URL)

		result = append(result, item)
	}

	return result, nil
}

// GetSyncData func
func (p *PostgreStorage) GetSyncData(id int64) ([]model.URLItem, error) {

	result := make([]model.URLItem, 0, 50)

	rows, err := p.DB.Query("select id, key, url from shortlinks where id > $1", id)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		item := model.URLItem{}
		rows.Scan(&item.ID, &item.Key, &item.URL)

		result = append(result, item)
	}

	return result, nil
}

// Save func
func (p *PostgreStorage) Save(item model.URLItem) error {

	_, err := p.DB.Exec("insert into shortlinks (id, key, url) values($1, $2, $3)on conflict on constraint pk_shortlinks do update set url=$3", item.ID, item.Key, item.URL)
	if err != nil {
		return err
	}

	return nil
}

// Delete func
func (p *PostgreStorage) Delete(id int64) error {

	_, err := p.DB.Exec("delete from shortlinks WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
