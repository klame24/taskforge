package main

import (
	"log"
	"net/http"
	"time"

	"taskforge/internal/auth"
	"taskforge/internal/config"
	"taskforge/internal/db"
	"taskforge/internal/handlers"
	"taskforge/internal/repository"
	"taskforge/internal/router"

	// "taskforge/internal/server"
	"taskforge/internal/service"
)

func main() {
	cfg := config.Load()

	jwtManager, err := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiration)
	if err != nil {
		log.Fatalf("Failed to create JWT manager: %v", err)
	}

	db, err := db.ConnectDB(cfg.DB_DSN)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtManager)
	userHandler := handlers.NewUserHandler(userService)

	appRouter := router.NewRouter(userHandler)
	handler := appRouter.SetupRoutes()

	log.Printf("Starting TaskForge API on :%s", cfg.HTTPPort)

	httpServer := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      handler, //
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
