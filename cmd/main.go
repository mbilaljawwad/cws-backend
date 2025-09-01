package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func service() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	return router
}

func main() {
	// TODO: Define proper routes and handlers.

	// create server context
	server := &http.Server{
		Addr:    ":8080",
		Handler: service(),
	}

	serverCtx, serverCancelCtx := context.WithCancel(context.Background())

	// signal for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-signalCh
		defer serverCancelCtx()

		shutdownCtx, cancelShutdownCtx := context.WithTimeout(serverCtx, time.Second*30)
		defer cancelShutdownCtx()
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out. Forcing exit.")
			}
		}()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("failed to shutdown server: %v", err)
		}
	}()

	log.Println("server started on port 8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %v", err)
	}

	// Wait for server context to be stopped.
	<-serverCtx.Done()
}
