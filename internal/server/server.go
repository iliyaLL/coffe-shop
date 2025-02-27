package server

import (
	"database/sql"
	"frappuccino/internal/handlers"
	"frappuccino/internal/service"
	"log/slog"
	"net/http"
	"os"
)

type server struct {
	port   string
	db     *sql.DB
	logger *slog.Logger
}

func NewServer(port string, db *sql.DB, logger *slog.Logger) *server {
	return &server{
		port:   port,
		db:     db,
		logger: logger,
	}
}

func (s *server) RunServer() {
	app := handlers.NewApplication(service.NewInventoryService(s.db))

	srv := &http.Server{
		Addr:     s.port,
		Handler:  app.Routes(),
		ErrorLog: slog.NewLogLogger(s.logger.Handler(), slog.LevelError),
	}

	s.logger.Info("starting server", "addr", srv.Addr)

	err := srv.ListenAndServe()
	s.logger.Error(err.Error())
	os.Exit(1)
}
