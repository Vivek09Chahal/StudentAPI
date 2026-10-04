package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vivek09Chahal/student_api/internal/config"
	"github.com/Vivek09Chahal/student_api/internal/http/handelers/student"
	"github.com/Vivek09Chahal/student_api/internal/storage/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// load config

	cfg := config.MustLoad()

	// database setup
	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage initialize", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /api/students", student.New(storage))

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("started server", slog.String("addr", server.Addr))
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("failed to start server: %v", err)
		}
	}()

	<-done

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("value", err.Error()))
	}

	slog.Info("server down successfully")

}
