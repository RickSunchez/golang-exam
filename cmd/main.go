package main

import (
	"last_lesson/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	ListenAndServe()
}

func ListenAndServe() {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.HandleConnection)
	http.ListenAndServe("127.0.0.1:8585", router)
}
