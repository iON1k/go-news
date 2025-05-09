package rss

import (
	"GoNews/pkg/models"
	"GoNews/pkg/storage"
	"encoding/xml"
	"log"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
)

// Обработчик загрузки RSS данных
type Fetcher struct {
	loaders []Loader
	store   storage.Store
}

// Конструктор обработчика загрузки RSS данных
func NewFetcher(loaders []Loader, store storage.Store) *Fetcher {
	return &Fetcher{loaders: loaders, store: store}
}

// Запуск обработчика загрузки RSS данных
// Загружает данные из источников с периодичностью syncPeriod, и склыдывает их в БД.
func (f Fetcher) Start(syncPeriod int) {
	chNews := make(chan []models.FullNews)
	chErrs := make(chan error)

	for _, loader := range f.loaders {
		go syncLoader(loader, chNews, chErrs, syncPeriod)
	}

	go func() {
		for news := range chNews {
			err := f.store.AddNews(news)
			if err != nil {
				chErrs <- err
			}
		}
	}()

	go func() {
		for err := range chErrs {
			log.Println("Feed fetching error:", err)
		}
	}()
}

func syncLoader(loader Loader, chNews chan<- []models.FullNews, chErrs chan<- error, syncPeriod int) {
	for {
		feed, err := fetchFeed(loader)
		if err != nil {
			chErrs <- err
		} else {
			chNews <- parseFeed(feed)
		}

		time.Sleep(time.Minute * time.Duration(syncPeriod))
	}
}

func fetchFeed(loader Loader) (Feed, error) {
	body, err := loader.LoadFeed()
	if err != nil {
		return Feed{}, err
	}
	var feed Feed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return Feed{}, err
	}

	return feed, nil
}

func parseFeed(feed Feed) []models.FullNews {
	var result []models.FullNews
	for _, item := range feed.Chanel.Items {
		news := models.FullNews{
			Title:   item.Title,
			Content: strip.StripTags(item.Description),
			Link:    item.Link,
			PubTime: parsePubTime(item.PubDate),
		}
		result = append(result, news)
	}
	return result
}

func parsePubTime(date string) int64 {
	result, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", date)
	if err != nil {
		result, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", date)
	}
	if err != nil {
		return 0
	}

	return result.Unix()
}
