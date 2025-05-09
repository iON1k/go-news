package memdb

import (
	"GoNews/pkg/models"
	"errors"
	"strings"
	"sync"
)

// Хранилище данных в памяти
type Store struct {
	news []models.FullNews
	mut  sync.Mutex
}

// Конструктор хранилища.
func New() *Store {
	return &Store{make([]models.FullNews, 0), sync.Mutex{}}
}

// Получение публикаций по фильтру
func (s *Store) NewsList(titleFilter string, from int64, to int64, offset int, count int) ([]models.ShortNews, error) {
	s.mut.Lock()
	defer s.mut.Unlock()

	result := []models.ShortNews{}
	curOffset := 0
	for _, news := range s.news {
		if news.PubTime < from || to != 0 && news.PubTime > to {
			continue
		}

		if !strings.Contains(strings.ToLower(news.Title), strings.ToLower(titleFilter)) {
			continue
		}

		if curOffset >= offset {
			shortNews := models.ShortNews{ID: news.ID, Title: news.Title, PubTime: news.PubTime, Link: news.Link}
			result = append(result, shortNews)

			if count != 0 && len(result) == count {
				return result, nil
			}
		}

		curOffset++
	}

	return result, nil
}

// Получение деталей публикации
func (s *Store) NewsDetails(id int) (models.FullNews, error) {
	for _, news := range s.news {
		if news.ID == id {
			return news, nil
		}
	}

	return models.FullNews{}, errors.New("NEWS NOT FOUND")
}

// Добавление новых публикаций
func (s *Store) AddNews(news []models.FullNews) error {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.news = append(s.news, news...)
	return nil
}
