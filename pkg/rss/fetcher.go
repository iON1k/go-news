package rss

import (
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
	chPosts := make(chan []storage.Post)
	chErrs := make(chan error)

	for _, loader := range f.loaders {
		go syncLoader(loader, chPosts, chErrs, syncPeriod)
	}

	go func() {
		for posts := range chPosts {
			err := f.store.AddPosts(posts)
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

func syncLoader(loader Loader, chPosts chan<- []storage.Post, chErrs chan<- error, syncPeriod int) {
	for {
		feed, err := fetchFeed(loader)
		if err != nil {
			chErrs <- err
		} else {
			chPosts <- parseFeed(feed)
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

func parseFeed(feed Feed) []storage.Post {
	var result []storage.Post
	for _, item := range feed.Chanel.Items {
		post := storage.Post{
			Title:   item.Title,
			Content: strip.StripTags(item.Description),
			Link:    item.Link,
			PubTime: parsePubTime(item.PubDate),
		}
		result = append(result, post)
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
