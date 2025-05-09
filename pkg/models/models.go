package models

// Короткая модель публикации
type ShortNews struct {
	ID      int    `json:"id"`       // идентификатор публикации
	Title   string `json:"title"`    // заголовок публикации
	PubTime int64  `json:"pub_time"` // время публикации
	Link    string `json:"link"`     // ссылка на источник
}

// Полная модель публикации
type FullNews struct {
	ID      int    `json:"id"`       // идентификатор публикации
	Title   string `json:"title"`    // заголовок публикации
	Content string `json:"content"`  // содержание публикации
	PubTime int64  `json:"pub_time"` // время публикации
	Link    string `json:"link"`     // ссылка на источник
}
