package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/memdb"
	"net/http"
)

func main() {
	// db, err := postgres.New("postgres://TEST_DB")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	db := memdb.New()
	db.AddPost(storage.Post{Title: "Test 1", Content: "Test content 1", Link: "https://google.com"})
	db.AddPost(storage.Post{Title: "Test 2", Content: "Test content 2", Link: "https://google.com"})
	api := api.New(db)
	http.ListenAndServe(":80", api.Router())
}
