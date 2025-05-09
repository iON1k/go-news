package storage

import "GoNews/pkg/models"

// Хранилище данных
type Store interface {
	NewsList(titleFilter string, from int64, to int64, offset int, count int) ([]models.ShortNews, error) // Получение публикаций по фильтру
	NewsDetails(id int) (models.FullNews, error)                                                          // Получение деталей публикации
	AddNews(news []models.FullNews) error                                                                 // Добавление новых публикаций
}
