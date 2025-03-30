package storage

// Доменная модель публикации
type Post struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

// Хранилище данных
type Store interface {
	Posts(n int) ([]Post, error) // Получение последних n публикаций
	AddPost(p Post) error        // Добавление новой публикации
}
