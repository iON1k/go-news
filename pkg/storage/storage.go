package storage

import "news/pkg/models"

// Хранилище данных
type Store interface {
	NewsPage(titleFilter string, from int64, to int64, page int) (models.NewsPage, error) // Получение страницы с публикациями
	NewsDetails(id int) (models.FullNews, error)                                          // Получение деталей публикации
	AddNews(news []models.FullNews) error                                                 // Добавление новых публикаций
}
