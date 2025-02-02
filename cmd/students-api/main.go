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

	"github.com/agodse21/students-go-api/internal/config"
	"github.com/agodse21/students-go-api/internal/http/handlers/student"
	"github.com/agodse21/students-go-api/internal/storage/sqlite"
)

func main() {
	// Load config
	cfg := config.MustLoad()
	// db setup

	storage, er := sqlite.New(cfg)

	if er != nil {
		log.Fatal(er)
	}

	slog.Info("Connected to database", slog.String("env", cfg.Env))
	// setup router

	router := http.NewServeMux()

	router.HandleFunc("POST /students/create", student.New(storage))
	router.HandleFunc("GET /students/{id}", student.GetById(storage))
	router.HandleFunc("GET /students", student.GetAll(storage))

	// setup server

	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}
	// start server

	slog.Info("Starting server", slog.String("address", cfg.Address))
	// fmt.Printf("Server is running %s", cfg.HttpServer.Address)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Graceful shutdown
	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	<-done

	slog.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server gracefully shutdown")

}

// run command
// go run cmd/students-api/main.go  -config config/local.yaml
