package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ShiftyX1/Facade/internal/config"
)


type Server struct {
	httpServer *http.Server
	config     *config.Config
}


func New(cfg *config.Config, handler http.Handler) *Server {
	
	var h http.Handler = handler

	
	h = RecoveryMiddleware(h)

	
	if cfg.Global.LogRequests {
		h = LoggingMiddleware(h)
	}

	
	if cfg.Global.CORS {
		h = CORSMiddleware(h)
	}

	
	if cfg.Global.Delay > 0 {
		h = DelayMiddleware(cfg.Global.Delay)(h)
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      h,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		config:     cfg,
	}
}


func (s *Server) Start() error {
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	
	go func() {
		slog.Info("Starting server", "addr", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	
	<-quit
	slog.Info("Shutting down server...")

	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	
	if err := s.httpServer.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		return err
	}

	slog.Info("Server shutdown complete")
	return nil
}