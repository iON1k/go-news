package postgres

import (
	"GoNews/pkg/models"
	"context"
	"math"

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

// Получение публикаций по фильтру
func (s *Store) NewsList(titleFilter string, from int64, to int64, offset int, count int) ([]models.ShortNews, error) {
	if count <= 0 {
		count = 1000
	}

	if to <= 0 {
		to = math.MaxInt64
	}

	titleFilter = "%" + titleFilter + "%"

	rows, err := s.db.Query(
		context.Background(),
		`
		SELECT news.id AS id, title, pub_time, link
		FROM news
		WHERE pub_time >= $1 AND pub_time <= $2 AND title ILIKE $3
		ORDER BY pub_time DESC
		OFFSET $4
		LIMIT $5;
		`,
		from,
		to,
		titleFilter,
		offset,
		count,
	)

	if err != nil {
		return nil, err
	}

	result := []models.ShortNews{}
	for rows.Next() {
		var p models.ShortNews
		err := rows.Scan(
			&p.ID,
			&p.Title,
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

func (s *Store) NewsDetails(id int) (models.FullNews, error) {
	var news models.FullNews

	err := s.db.QueryRow(
		context.Background(),
		`
		SELECT news.id AS id, title, content, pub_time, link
		FROM news
		WHERE id = $1;
		`,
		id,
	).Scan(
		&news.ID,
		&news.Title,
		&news.Content,
		&news.PubTime,
		&news.Link,
	)

	return news, err
}

// Добавление новых публикаций
func (s *Store) AddNews(news []models.FullNews) error {
	ctx := context.Background()
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	for _, n := range news {
		_, err := tx.Exec(
			context.Background(),
			`
			INSERT INTO news (title, content, pub_time, link) 
			VALUES ($1, $2, $3, $4)
			ON CONFLICT DO NOTHING;
			`,
			n.Title,
			n.Content,
			n.PubTime,
			n.Link,
		)

		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	return tx.Commit(ctx)
}
