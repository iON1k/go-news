package postgres

import (
	"context"
	"news/pkg/models"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const TEST_NEWS_PAGE_SIZE = 3

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

func TestNewsPageWithoutFilters(t *testing.T) {
	ctx := setup(t)
	defer ctx.teardown()

	got_page, err := ctx.s.NewsPage("", 0, 0, 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(got_page.News) != 2 || got_page.News[0].Title != "Test2" || got_page.News[1].Title != "Test1" {
		t.Fatal("Got wrong news from DB")
	}

	if got_page.Paging.Index != 1 || got_page.Paging.Count != 2 || got_page.Paging.Size != TEST_NEWS_PAGE_SIZE {
		t.Fatal("Got wrong paging")
	}
}

func TestNewsPageWithPubTimeFilter(t *testing.T) {
	ctx := setup(t)
	defer ctx.teardown()

	got_page, err := ctx.s.NewsPage("", 2, 3, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(got_page.News) != 2 || got_page.News[0].Title != "Test4" || got_page.News[1].Title != "Test3" {
		t.Fatal("Got wrong news from DB")
	}

	if got_page.Paging.Index != 0 || got_page.Paging.Count != 1 || got_page.Paging.Size != TEST_NEWS_PAGE_SIZE {
		t.Fatal("Got wrong paging")
	}
}

func TestNewPageWithTitleFilter(t *testing.T) {
	ctx := setup(t)
	defer ctx.teardown()

	got_page, err := ctx.s.NewsPage("ST1", 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(got_page.News) != 1 || got_page.News[0].Title != "Test1" {
		t.Fatal("Got wrong news from DB")
	}

	if got_page.Paging.Index != 0 || got_page.Paging.Count != 1 || got_page.Paging.Size != TEST_NEWS_PAGE_SIZE {
		t.Fatal("Got wrong paging")
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

	return NewFromPGX(db, TEST_NEWS_PAGE_SIZE)
}
