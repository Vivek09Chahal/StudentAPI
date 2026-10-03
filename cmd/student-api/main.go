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
)

func main() {
    // load config
 
    cfg := config.MustLoad()
    
    // database setup
    
    // setup router
    router := http.NewServeMux()

    router.HandleFunc("GET /api/students", student.New())
    
    // setup server
    server := http.Server {
        Addr: cfg.Addr,
        Handler: router,
    }

    done := make(chan os.Signal, 1)

    signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        err := server.ListenAndServe()
        if err != nil {
            log.Fatal("failed to start server")
        }
    }()

    <-done  

    slog.Info("shutting down server")
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    err := server.Shutdown(ctx)
    if err != nil {
        slog.Error("failed to shutdown server", slog.String("value", err.Error()))
    }

    slog.Info("server down successfully")
    
}