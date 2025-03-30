package postgres

import (
	"GoNews/pkg/storage"
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
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
		t.Error(err)
	}

	if len(posts) != 3 {
		t.Fatalf("Wrong posts count")
	}
}

func TestStorage(t *testing.T) {
	ctx := setup(t)
	defer ctx.teardown()

	post1 := storage.Post{Title: "Test1", PubTime: 0}
	post2 := storage.Post{Title: "Test2", PubTime: 1}
	post3 := storage.Post{Title: "Test3", PubTime: 2}
	post4 := storage.Post{Title: "Test4", PubTime: 3}
	post5 := storage.Post{Title: "Test5", PubTime: 4}
	posts_to_add := []storage.Post{post3, post1, post2, post5, post4}

	for _, post := range posts_to_add {
		err := ctx.s.AddPost(post)
		if err != nil {
			t.Error(err)
		}
	}

	got, err := ctx.s.Posts(3)
	if err != nil {
		t.Error(err)
	}

	if len(got) != 3 || got[0].Title != "Test5" || got[2].Title != "Test3" {
		t.Fatalf("Got wrong posts from DB")
	}
}

func makeStorage(t *testing.T) *Store {
	db, err := pgxpool.New(context.Background(), "postgres://postgres:post_irony88@104.252.127.170:5432")
	if err != nil {
		t.Error(err)
	}

	bytes, err := os.ReadFile("schema.sql")
	if err != nil {
		t.Error(err)
		db.Close()
		return nil
	}

	_, err = db.Exec(context.Background(), string(bytes))
	if err != nil {
		t.Error(err)
		db.Close()
		return nil
	}

	return NewFromPGX(db)
}
