package main

import (
	"encoding/json"
	"log"
	"net/http"
	"news/pkg/api"
	"news/pkg/rss"
	"news/pkg/storage/postgres"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем файл окружения
	godotenv.Load()

	db_conn := os.Getenv("DB")
	if db_conn == "" {
		log.Fatal("No environment for DB")
	}

	// Читаем файл конфигурации
	c, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config config
	err = json.Unmarshal(c, &config)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем подключение к БД
	store, err := postgres.New(db_conn)
	if err != nil {
		log.Fatal(err)
	}

	defer store.Close()

	// Создаем загрузчики данных RSS
	var rss_loaders []rss.Loader
	for _, url := range config.Urls {
		l := rss.NewHttpLoader(url)
		rss_loaders = append(rss_loaders, l)
	}

	// Запускаем обработчик данных RSS
	fetcher := rss.NewFetcher(rss_loaders, store)
	fetcher.Start(config.SyncPeriod)

	// Запускаем API
	api := api.New(store)
	log.Print("Server is starting...")
	http.ListenAndServe(":8080", api.Router())
	log.Print("Server has been stopped.")
}
