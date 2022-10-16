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
	Delete(key string) error
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

	_, err := p.DB.Exec("create table if not exists data(key varchar, url varchar, constraint pk_data primary key(key))")
	if err != nil {
		return err
	}

	return nil
}

// Get func
func (p *PostgreStorage) Get() ([]model.URLItem, error) {

	result := make([]model.URLItem, 0, 50)

	rows, err := p.DB.Query("select key, url from data")
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		item := model.URLItem{}
		rows.Scan(&item.Key, &item.URL)

		result = append(result, item)
	}

	return result, nil
}

// Save func
func (p *PostgreStorage) Save(item model.URLItem) error {

	_, err := p.DB.Exec("insert into data (key, url) values($1, $2)on conflict on constraint pk_data do update set url=$2", item.Key, item.URL)
	if err != nil {
		return err
	}

	return nil
}

// Delete func
func (p *PostgreStorage) Delete(key string) error {

	_, err := p.DB.Exec("delete from data WHERE key = $1", key)
	if err != nil {
		return err
	}

	return nil
}