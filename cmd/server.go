package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"news/pkg/api"
	"news/pkg/rss"
	"news/pkg/storage"
	"news/pkg/storage/postgres"
)

// Конфигурация приложения
type config struct {
	URLS   []string `json:"rss"`
	Period int      `json:"request_period"`
}

type server struct {
	db  storage.DBInterface
	api *api.API
}

func main() {
	// Читаем файл конфигурации приложения
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Раскодируем данные файла конфигурации
	var c config
	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Fatal(err)
	}

	var srv server

	pgPass := os.Getenv("pgPass")
	if pgPass == "" {
		os.Exit(1)
	}
	// conn - строка подключения к базе данных
	pgConn := "postgres://postgres:" + pgPass + "@localhost:5432/news"

	// Инициализируем базу данных Postgres
	db, err := postgres.New(pgConn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Инициализируем хранилище сервера приложения
	srv.db = db

	// Создаём объект API и регистрируем обработчики
	srv.api = api.New(srv.db)

	// Запускаем парсер новостей, для каждой RSS-ленты
	// в отдельной горутине
	chNews := make(chan []storage.Post)
	chErrors := make(chan error)
	for _, url := range c.URLS {
		go parse(url, chNews, chErrors, c.Period)
	}

	// Запись потока новостей в базу данных
	go func() {
		for posts := range chNews {
			srv.db.StoreNews(posts)
		}
	}()

	// Обработка потока ошибок
	go func() {
		for err := range chErrors {
			log.Println("Ошибка парсера RSS-лент:", err)
		}
	}()

	// Запуск веб-сервер на порту 80
	http.ListenAndServe(":80", srv.api.Router())
}

// Парсер RSS-лент, пишет новости и ошибки в соответствующие каналы posts и errors
func parse(url string, posts chan<- []storage.Post, errors chan<- error, period int) {
	for {
		news, err := rss.ParseFeed(url)
		if err != nil {
			errors <- err
		}
		posts <- news
		time.Sleep(time.Duration(period) * time.Minute)
	}
}
