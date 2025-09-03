package main

import (
	"blockchain/database"
	"blockchain/routes"
	"blockchain/schedule"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gofr.dev/pkg/gofr"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDB()
	database.Migrate()

	r := routes.InitRoute()
	cron := gofr.New()

	cron.AddCronJob("*/5 * * * * *", "Mine Block", func(ctx *gofr.Context) {
		schedule.AutoMineBlock()
	})

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: corsHandler.Handler(r),

		//read timeout
		ReadTimeout: 15 * time.Second,
		//write timeout
		WriteTimeout: 15 * time.Second,
		//idle timeout
		IdleTimeout: 60 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	sqlDB, errDB := database.DB.DB()
	if errDB == nil {
		if errClose := sqlDB.Close(); errClose != nil {
			log.Fatal("Error closing database connection")
		} else {
			log.Fatal("Database connection closed")
		}
	} else {
		log.Fatal("Error closing database connection")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown failed:", err)
	}
}
