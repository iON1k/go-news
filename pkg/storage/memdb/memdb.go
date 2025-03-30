package memdb

import (
	"GoNews/pkg/storage"
	"sync"
)

// Хранилище данных в памяти
type Store struct {
	posts []storage.Post
	mut   sync.Mutex
}

// Конструктор хранилища.
func New() *Store {
	return &Store{make([]storage.Post, 0), sync.Mutex{}}
}

// Получение последних n публикаций
func (s *Store) Posts(n int) ([]storage.Post, error) {
	s.mut.Lock()
	defer s.mut.Unlock()
	var res_len int
	if n <= 0 {
		res_len = len(s.posts)
	} else {
		res_len = min(len(s.posts), n)
	}

	res := make([]storage.Post, res_len)
	copy(res, s.posts)

	return res, nil
}

// Добавление новой публикации
func (s *Store) AddPost(p storage.Post) error {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.posts = append(s.posts, p)
	return nil
}
