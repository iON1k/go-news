package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"news/pkg/models"
	"news/pkg/storage/memdb"
	"testing"
)

type TestContext struct {
	api *API
}

func setup(_ *testing.T) TestContext {
	db := memdb.New()
	db.AddNews([]models.FullNews{
		{ID: 1, Title: "Test1", Content: "Content1"},
		{ID: 2, Title: "Test2", Content: "Content2"},
		{ID: 3, Title: "Test3", Content: "Content3"},
	})
	api := New(db)

	return TestContext{api}
}

func TestNewsList(t *testing.T) {
	ctx := setup(t)

	req := httptest.NewRequest(http.MethodGet, "/news", nil)
	resp := httptest.NewRecorder()
	ctx.api.router.ServeHTTP(resp, req)

	if !(resp.Code == http.StatusOK) {
		t.Fatal("Wrong status code")
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var n_page models.NewsPage
	err = json.Unmarshal(b, &n_page)
	if err != nil {
		t.Fatal(err)
	}

	if len(n_page.News) != 3 || n_page.News[0].Title != "Test1" || n_page.News[1].Title != "Test2" ||
		n_page.News[2].Title != "Test3" || n_page.Paging.Count != 1 {
		t.Fatalf("Wrong data in response")
	}
}

func TestNewsDetails(t *testing.T) {
	ctx := setup(t)

	req := httptest.NewRequest(http.MethodGet, "/news/1", nil)
	resp := httptest.NewRecorder()
	ctx.api.router.ServeHTTP(resp, req)

	if !(resp.Code == http.StatusOK) {
		t.Fatal("Wrong status code")
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var n models.FullNews
	err = json.Unmarshal(b, &n)
	if err != nil {
		t.Fatal(err)
	}

	if n.ID != 1 || n.Title != "Test1" || n.Content != "Content1" {
		t.Fatalf("Wrong data in response")
	}
}
