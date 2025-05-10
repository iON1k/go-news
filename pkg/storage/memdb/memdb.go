package memdb

import (
	"errors"
	"math"
	"news/pkg/models"
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

// Получение страницы с публикациями
func (s *Store) NewsPage(titleFilter string, from int64, to int64, page int) (models.NewsPage, error) {
	paging := models.Paging{Index: page, Count: 1, Size: math.MaxInt}
	if page != 0 {
		return models.NewsPage{News: []models.ShortNews{}, Paging: paging}, nil
	}

	s.mut.Lock()
	defer s.mut.Unlock()
	s_news := []models.ShortNews{}
	for _, news := range s.news {
		sn := models.ShortNews{ID: news.ID, Title: news.Title, PubTime: news.PubTime, Link: news.Link}
		s_news = append(s_news, sn)
	}

	return models.NewsPage{News: s_news, Paging: paging}, nil
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
