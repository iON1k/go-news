package postgres

import (
	"GoNews/pkg/storage"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

// Конструктор хранилища с URL для коннекта к БД
func New(conn string) (*Store, error) {
	db, err := pgxpool.New(context.Background(), conn)

	if err != nil {
		return nil, err
	}

	return NewFromPGX(db), nil
}

// Конструктор хранилища с готовым коннектом к БД
func NewFromPGX(db *pgxpool.Pool) *Store {
	return &Store{db}
}

func (s *Store) Close() {
	s.db.Close()
}

// Получение последних n публикаций
func (s *Store) Posts(n int) ([]storage.Post, error) {
	if n <= 0 {
		n = 10
	}

	rows, err := s.db.Query(
		context.Background(),
		`
		SELECT posts.id AS id, title, content, pub_time, link
		FROM posts
		ORDER BY pub_time DESC
		LIMIT $1;
		`,
		n,
	)

	if err != nil {
		return nil, err
	}

	var result []storage.Post
	for rows.Next() {
		var p storage.Post
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)

		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	return result, rows.Err()
}

// Добавление новой публикации
func (s *Store) AddPost(p storage.Post) error {
	_, err := s.db.Exec(
		context.Background(),
		`
		INSERT INTO posts (title, content, pub_time, link) 
		VALUES ($1, $2, $3, $4);
		`,
		p.Title,
		p.Content,
		p.PubTime,
		p.Link,
	)

	return err
}
