package postgres

import (
	"context"
	"math"
	"news/pkg/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

const NEWS_PAGE_SIZE = 10

// Хранилище данных.
type Store struct {
	db     *pgxpool.Pool
	p_size int
}

// Конструктор хранилища с URL для коннекта к БД
func New(conn string) (*Store, error) {
	db, err := pgxpool.New(context.Background(), conn)

	if err != nil {
		return nil, err
	}

	return NewFromPGX(db, NEWS_PAGE_SIZE), nil
}

// Конструктор хранилища с готовым коннектом к БД
func NewFromPGX(db *pgxpool.Pool, p_size int) *Store {
	return &Store{db, p_size}
}

func (s *Store) Close() {
	s.db.Close()
}

// Получение страницы с публикациями
func (s *Store) NewsPage(titleFilter string, from int64, to int64, page int) (models.NewsPage, error) {
	if to <= 0 {
		to = math.MaxInt64
	}

	titleFilter = "%" + titleFilter + "%"

	var allRowsCount int
	err := s.db.QueryRow(
		context.Background(),
		`
		SELECT COUNT(*)
		FROM news
		WHERE pub_time >= $1 AND pub_time <= $2 AND title ILIKE $3;
		`,
		from,
		to,
		titleFilter,
	).Scan(&allRowsCount)

	if err != nil {
		return models.NewsPage{}, err
	}

	pageRows, err := s.db.Query(
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
		page*s.p_size,
		s.p_size,
	)

	if err != nil {
		return models.NewsPage{}, err
	}

	page_news := []models.ShortNews{}
	for pageRows.Next() {
		var n models.ShortNews
		err := pageRows.Scan(
			&n.ID,
			&n.Title,
			&n.PubTime,
			&n.Link,
		)

		if err != nil {
			return models.NewsPage{}, err
		}
		page_news = append(page_news, n)
	}

	pagesCount := allRowsCount / s.p_size
	if allRowsCount%s.p_size != 0 {
		pagesCount++
	}

	n_page := models.NewsPage{
		News:   page_news,
		Paging: models.Paging{Index: page, Count: pagesCount, Size: s.p_size},
	}

	return n_page, pageRows.Err()
}

func (s *Store) NewsDetails(id int) (models.FullNews, error) {
	var n models.FullNews

	err := s.db.QueryRow(
		context.Background(),
		`
		SELECT news.id AS id, title, content, pub_time, link
		FROM news
		WHERE id = $1;
		`,
		id,
	).Scan(
		&n.ID,
		&n.Title,
		&n.Content,
		&n.PubTime,
		&n.Link,
	)

	return n, err
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
