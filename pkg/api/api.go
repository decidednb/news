package api

import (
	"encoding/json"
	"net/http"
	"news/pkg/storage"
	"strconv"

	"github.com/gorilla/mux"
)

// Программный интерфейс сервера
type API struct {
	db     storage.DBInterface
	router *mux.Router
}

// Конструктор API
func New(db storage.DBInterface) *API {
	api := API{
		db: db,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// Получение маршрутизатора запросов.
func (api *API) Router() *mux.Router {
	return api.router
}

// Обработчики API.
func (api *API) endpoints() {
	// Получение последних n новостей
	api.router.HandleFunc("/news/{n}", api.lastNewsHandler).Methods(http.MethodGet, http.MethodOptions)
	// Веб-приложение
	api.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

func (api *API) lastNewsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	s := mux.Vars(r)["n"]
	n, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	news, err := api.db.LastNews(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(news)
}
