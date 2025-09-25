package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"tasklist/db"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"tasklist/internal/handler"
	"tasklist/internal/repository"
	"tasklist/internal/service"
	"tasklist/pkg/config"
	"tasklist/pkg/logger"
)

// @title TodoList API
// @version 1.0
// @description This is a TaskList server built with Go, Gin and PostgreSQL.
//
// @host localhost:1232
// @BasePath /api
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	log := logger.New(cfg.LogLevel)
	log.Info("Starting server...")

	db, err := db.InitDB(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer db.Pool.Close()

	log.Info("Connected to database...")

	repo := repository.NewRepository(db)
	service := service.NewService(repo, cfg)
	handler := handler.NewHandler(service, log)
	r := handler.Init()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler: r,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go startOverdueChecker(ctx, repo.TaskRepo, log)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("shutting down server...")

	ctxShut, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShut); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Info("server exiting")
}

func startOverdueChecker(ctx context.Context, repo repository.Task, log *logrus.Logger) {
	ticker := time.NewTicker(time.Hour * 1)
	defer ticker.Stop()

	checkOnce := func() {
		log.Info("checking overdue tasks...")
		tasks, err := repo.ListAll(ctx)
		if err != nil {
			log.WithError(err).Error("list tasks failed")
			return
		}
		now := time.Now()
		for _, task := range tasks {
			if task.Deadline != nil && task.Deadline.Before(now) {
				log.WithFields(logrus.Fields{
					"id":       task.ID,
					"deadline": task.Deadline.Format(time.RFC3339),
				}).Info("task is overdue")
				if err := repo.MarkOverdued(ctx, task.ID); err != nil {
					log.WithError(err).WithField("id", task.ID).Error("mark overdue failed")
				}
			}
		}
		log.Info("overdue check complete")
	}

	checkOnce()
	for {
		select {
		case <-ctx.Done():
			log.Info("overdue checker stopping")
			return
		case <-ticker.C:
			checkOnce()
		}
	}
}
