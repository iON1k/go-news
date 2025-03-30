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
	api.router.HandleFunc("/news/{n}", api.news).Methods(http.MethodGet, http.MethodOptions)
	api.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

func (api *API) news(w http.ResponseWriter, r *http.Request) {
	n_param := mux.Vars(r)["n"]
	n, err := strconv.Atoi(n_param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	posts, err := api.store.Posts(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(posts)
}
