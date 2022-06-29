package main

import (
	"last_lesson/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

// REPORT:
// getRandomIntBetweenValues - не обрабатывает верхнюю границу (возвращал всегда status closed для accendent)

func main() {
	ListenAndServe()
}

func ListenAndServe() {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.HandleConnection)
	http.ListenAndServe("127.0.0.1:8585", router)
}
