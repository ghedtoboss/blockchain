package main

import (
	"blockchain/controllers"
	"blockchain/database"
	"blockchain/models"
	"blockchain/routes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDB()
	database.Migrate()

	r := routes.InitRoute()

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

	
		bc := controllers.Blockchain{}
		// 1. Genesis Block
		genesis := bc.CreateGenesisBlock(time.Now().Unix())
		fmt.Println("Genesis Hash: ", genesis.Hash)

		// 2. Block 1
		block1 := bc.CreateBlock(genesis, []models.Transaction{
			{
				From:      "Veli",
				To:        "Ayşe",
				Currency:  "BTC",
				Amount:    10,
				Fee:       0.1,
				Signature: "",
			},
		}, time.Now().Unix())
		fmt.Println("Block 1 Hash: ", block1.Hash)

		// 3. Block 2
		block2 := bc.CreateBlock(block1, []models.Transaction{
			{
				From:      "Ayşe",
				To:        "Veli",
				Currency:  "BTC",
				Amount:    5,
				Fee:       0.05,
				Signature: "",
			},
		}, time.Now().Unix())
		fmt.Println("Block 2 Hash: ", block2.Hash)

		// Validate block
		if !bc.ValidateBlock(block1, block2) {
			fmt.Println("Block 1 is invalid")
			return
		} else {
			fmt.Println("Block 1 is valid")
		}

		// Validate chain
		if bc.ValidateChain() {
			fmt.Println("Blockchain is valid")
		} else {
			fmt.Println("Blockchain is invalid")
		}

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
