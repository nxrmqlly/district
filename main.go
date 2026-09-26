package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("hello world")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"bl\": \"hello, world\"}"))
	})

	srv := http.Server{
		Addr:    ":2468",
		Handler: mux,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println(err)
		os.Exit(1)
	}

}
