package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Run(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world from server!")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	}()

	log.Println("Starting server on :8080")
	err := srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil // expected after Shutdown
	}
	return err
}
