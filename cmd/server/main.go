package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bot/internal/api"
	"bot/internal/repository"
	"bot/internal/service"
	"bot/pkg/logger"
)

func main() {

	l := logger.NewLogger()
	l.Info("Starting voting system server...")


	tarantoolAddr := getEnv("TARANTOOL_ADDR", "tarantool:3301")
	serverAddr := getEnv("SERVER_ADDR", ":8080")


	repo, err := repository.NewTarantoolRepository(tarantoolAddr, l)
	if err != nil {
		l.Fatal("Failed to connect to Tarantool", err)
	}
	defer repo.Close()

	votingService := service.NewVotingService(repo, l)

	
	router := api.SetupRoutes(votingService, l)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	
	go func() {
		l.Info("Server listening on " + serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal("Error starting server", err)
		}
	}()


	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	l.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		l.Fatal("Server forced to shutdown", err)
	}

	l.Info("Server exited properly")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
