package postgres

import (
	"GoNews/pkg/storage"
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
	return TestContext{makeStorage(t)}
}

func (c TestContext) teardown() {
	if c.s != nil {
		c.s.Close()
	}
}

func TestPosts(t *testing.T) {
	c := setup(t)
	defer c.teardown()

	posts, err := c.s.Posts(5)
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 3 {
		t.Fatal("Wrong posts count")
	}
}

func TestStorage(t *testing.T) {
	ctx := setup(t)
	defer ctx.teardown()

	post1 := storage.Post{Title: "Test1", PubTime: 0, Link: "Link1"}
	post2 := storage.Post{Title: "Test2", PubTime: 1, Link: "Link2"}
	post3 := storage.Post{Title: "Test3", PubTime: 2, Link: "Link3"}
	post4 := storage.Post{Title: "Test4", PubTime: 3, Link: "Link4"}
	post5 := storage.Post{Title: "Test5", PubTime: 4, Link: "Link5"}
	posts_to_add := []storage.Post{post3, post1, post2, post5, post2, post3, post4}

	err := ctx.s.AddPosts(posts_to_add)
	if err != nil {
		t.Fatal(err)
	}

	got, err := ctx.s.Posts(3)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 3 || got[0].Title != "Test5" || got[2].Title != "Test3" {
		t.Fatal("Got wrong posts from DB")
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
