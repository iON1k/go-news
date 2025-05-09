package api

import (
	"GoNews/pkg/storage"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Программный интерфейс сервера
type API struct {
	store  storage.Store
	router *mux.Router
}

// Конструктор объекта API
func New(store storage.Store) *API {
	api := API{
		store: store,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// Маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.router
}

func (api *API) endpoints() {
	api.router.HandleFunc("/news", api.newsList).Methods(http.MethodGet)
	api.router.HandleFunc("/news/{id}", api.newsDetails).Methods(http.MethodGet)
}

func (api *API) newsList(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")

	fromStr := r.URL.Query().Get("from")
	from, _ := strconv.ParseInt(fromStr, 10, 64)

	toStr := r.URL.Query().Get("to")
	to, _ := strconv.ParseInt(toStr, 10, 64)

	offsetStr := r.URL.Query().Get("offset")
	offset, _ := strconv.Atoi(offsetStr)

	countStr := r.URL.Query().Get("count")
	count, _ := strconv.Atoi(countStr)

	news, err := api.store.NewsList(title, from, to, offset, count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(news)
}

func (api *API) newsDetails(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Id expected", http.StatusBadRequest)
		return
	}

	news, err := api.store.NewsDetails(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(news)
}
