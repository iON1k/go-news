package models

// Модель с информацией о пагинации
type Paging struct {
	Index int `json:"index"` // текущая страница
	Count int `json:"count"` // общее количество страниц
	Size  int `json:"size"`  // размер страницы
}

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

// Модель страницы со списком публикаций
type NewsPage struct {
	News   []ShortNews `json:"news"`   // список публикаций
	Paging Paging      `json:"paging"` // информация о пагинации
}
