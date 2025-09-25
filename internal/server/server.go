package server

import (
	"context"
	"cws-backend/internal/config"
	"cws-backend/internal/database"
	"cws-backend/internal/routes"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	ctx context.Context
	cfg *config.Config
	srv *http.Server
}

func NewServer(
	ctx context.Context,
	cfg *config.Config,
	dbm *database.DBManager,
) *Server {
	server := &Server{
		cfg: cfg,
		ctx: ctx,
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.AppPort),
			Handler: routes.NewRoutes(dbm).AppRouter,
		},
	}
	return server
}

func (s *Server) Start() {
	go func() {
		log.Println("server started on port", s.cfg.AppPort)
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}

		<-s.ctx.Done()
		log.Println("server has been stopped")
		if s.ctx.Err() == context.DeadlineExceeded {
			log.Fatal("graceful shutdown timed out. Forcing exit.")
		}

		if err := s.srv.Shutdown(s.ctx); err != nil {
			log.Fatalf("failed to shutdown server: %v", err)
		}
	}()
}
