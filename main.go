package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {

}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	// 1. Create a new ServeMux
	mux := http.NewServeMux()

	// 2. Create a new http.Server struct.
	// - Use the new "ServeMux" as the server's handler
	// - Set the .Addr field to ":8080"
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// 3. Register the homeHandler function to handle requests to the root path ("/")
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/healthz", http.HandlerFunc(healthzHandler))

	// 4. Start the server
	fmt.Println("Starting server on :8080")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
