package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"news/pkg/storage"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ctxRequestIdKey struct{}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

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
	api.router.Use(requestIdValidator)
	api.router.Use(requestLogger)

	api.router.Methods(http.MethodGet).Path("/news").HandlerFunc(api.newsList)
	api.router.Methods(http.MethodGet).Path("/news/{id}").HandlerFunc(api.newsDetails)
}

func (api *API) newsList(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("s")

	fromStr := r.URL.Query().Get("from")
	from, _ := strconv.ParseInt(fromStr, 10, 64)

	toStr := r.URL.Query().Get("to")
	to, _ := strconv.ParseInt(toStr, 10, 64)

	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)

	new_page, err := api.store.NewsPage(title, from, to, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(new_page)
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

func requestIdValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r_key := "request_id"
		url := r.URL
		req_id := url.Query().Get(r_key)
		if req_id == "" {
			req_id = uuid.New().String()
			addQueryToUrl(url, map[string]string{r_key: req_id})
		}

		ctx := context.WithValue(r.Context(), ctxRequestIdKey{}, req_id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req_id := getRequestId(r.Context())
		d_format := "2006-01-02 15:04:05"

		log.Printf(
			"REQUEST - ID: %v IP: %v TIME: %v",
			req_id,
			r.RemoteAddr,
			time.Now().Format(d_format),
		)

		status_w := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(status_w, r)

		log.Printf(
			"RESPONSE - ID: %v STATIS: %v TIME %v",
			req_id,
			status_w.status,
			time.Now().Format(d_format),
		)
	})
}

func addQueryToUrl(url *url.URL, params map[string]string) {
	q := url.Query()
	for k, v := range params {
		q.Add(k, v)
	}

	url.RawQuery = q.Encode()
}

func getRequestId(ctx context.Context) string {
	id, ok := ctx.Value(ctxRequestIdKey{}).(string)
	if !ok {
		return ""
	}
	return id
}
