package rss

import (
	"GoNews/pkg/storage/memdb"
	"testing"
	"time"
)

type TestContext struct {
	f *Fetcher
	s *memdb.Store
}

func setup() TestContext {
	loader := NewFileLoader("./test_rss.xml")
	loaders := []Loader{loader, loader, loader}
	store := memdb.New()

	return TestContext{
		NewFetcher(loaders, store),
		store,
	}
}

func TestFetching(t *testing.T) {
	ctx := setup()
	ctx.f.Start(1)

	time.Sleep(time.Second)
	posts, _ := ctx.s.Posts(0)

	if len(posts) != 6 {
		t.Fatal("Wrong fetched posts count")
	}
}
