package postgres

import (
	"GoNews/pkg/models"
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type TestContext struct {
	s *Store
}

func setup(t *testing.T) TestContext {
	ctx := TestContext{makeStorage(t)}

	news1 := models.FullNews{Title: "Test1", PubTime: 0, Link: "Link1"}
	news2 := models.FullNews{Title: "Test2", PubTime: 1, Link: "Link2"}
	news3 := models.FullNews{Title: "Test3", PubTime: 2, Link: "Link3"}
	news4 := models.FullNews{Title: "Test4", PubTime: 3, Link: "Link4"}
	news5 := models.FullNews{Title: "Test5", PubTime: 4, Link: "Link5"}
	news_to_add := []models.FullNews{news3, news1, news2, news5, news2, news3, news4}

	err := ctx.s.AddNews(news_to_add)
	if err != nil {
		t.Fatal(err)
	}

	return ctx
}

func (c TestContext) teardown() {
	if c.s != nil {
		c.s.Close()
	}
}

func TestNewListWithoutFilters(t *testing.T) {
	ctx := setup(t)
	defer ctx.teardown()

	got_news, err := ctx.s.NewsList("", 0, 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(got_news) != 5 || got_news[0].Title != "Test5" || got_news[4].Title != "Test1" {
		t.Fatal("Got wrong news from DB")
	}
}

func TestNewListWithCount(t *testing.T) {
	ctx := setup(t)
	defer ctx.teardown()

	got_news, err := ctx.s.NewsList("", 0, 0, 0, 3)
	if err != nil {
		t.Fatal(err)
	}

	if len(got_news) != 3 || got_news[0].Title != "Test5" || got_news[2].Title != "Test3" {
		t.Fatal("Got wrong news from DB")
	}
}

func TestNewListWithPubTimeFilter(t *testing.T) {
	ctx := setup(t)
	defer ctx.teardown()

	got_news, err := ctx.s.NewsList("", 2, 3, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(got_news) != 2 || got_news[0].Title != "Test4" || got_news[1].Title != "Test3" {
		t.Fatal("Got wrong news from DB")
	}
}

func TestNewListWithOffset(t *testing.T) {
	ctx := setup(t)
	defer ctx.teardown()

	got_news, err := ctx.s.NewsList("", 0, 0, 3, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(got_news) != 2 || got_news[0].Title != "Test2" || got_news[1].Title != "Test1" {
		t.Fatal("Got wrong news from DB")
	}
}

func TestNewListWithTitleFilter(t *testing.T) {
	ctx := setup(t)
	defer ctx.teardown()

	got_news, err := ctx.s.NewsList("ST1", 0, 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(got_news) != 1 || got_news[0].Title != "Test1" {
		t.Fatal("Got wrong news from DB")
	}
}

func TestNewsDetails(t *testing.T) {
	ctx := setup(t)
	defer ctx.teardown()

	got_news, err := ctx.s.NewsDetails(1)
	if err != nil {
		t.Fatal(err)
	}

	if got_news.Title != "Test3" {
		t.Fatal("Got wrong news from DB")
	}
}

func makeStorage(t *testing.T) *Store {
	err := godotenv.Load()
	if err != nil {
		t.Fatal(err)
	}

	db_conn := os.Getenv("TEST_DB")
	if db_conn == "" {
		t.Fatal("No environment for DB")
	}

	db, err := pgxpool.New(context.Background(), db_conn)
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := os.ReadFile("schema.sql")
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	_, err = db.Exec(context.Background(), string(bytes))
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	return NewFromPGX(db)
}
