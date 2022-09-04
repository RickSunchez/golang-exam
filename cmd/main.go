package main

import (
	"fmt"
	"last_lesson/internal/handlers"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

func main() {
	ListenAndServe()
}

func ListenAndServe() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		router := mux.NewRouter()
		router.HandleFunc("/", handlers.HandleConnection)
		http.ListenAndServe("127.0.0.1:8585", router)
	}()

	fmt.Println("Served at http://127.0.0.1:8585")

	<-done

	fmt.Println("\nServer stopped")
}
