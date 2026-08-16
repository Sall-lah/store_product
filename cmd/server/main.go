package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sall-lah/store_product/internal/cache"
	"github.com/Sall-lah/store_product/internal/config"
	"github.com/Sall-lah/store_product/internal/db"
	"github.com/Sall-lah/store_product/internal/handler"
	"github.com/Sall-lah/store_product/internal/repository"
	"github.com/Sall-lah/store_product/internal/service"
)

func main() {
	log.Println("[INFO] Starting store_product microservice...")

	// 1. Load application configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[FATAL] Failed to load configuration: %v", err)
	}

	// 2. Initialize Redis client (Local port 6739 by default)
	cacheClient := cache.NewClient(cfg)

	// 3. Initialize Prisma Client Go for Supabase PostgreSQL
	var prismaClient *db.PrismaClient
	if cfg.DatabaseURL != "" {
		prismaClient = db.NewClient(db.WithDatasourceURL(cfg.DatabaseURL))
	} else {
		prismaClient = db.NewClient()
	}

	if err := prismaClient.Prisma.Connect(); err != nil {
		log.Printf("[WARN] Initial database connection error: %v. Running in disconnected standby mode until valid DATABASE_URL is supplied.", err)
	} else {
		log.Println("[INFO] Connected successfully to Supabase PostgreSQL database.")
	}

	defer func() {
		if err := prismaClient.Prisma.Disconnect(); err != nil {
			log.Printf("[WARN] Error disconnecting from Prisma: %v", err)
		}
	}()

	// 4. Wire repository, service, and transport layers (Dependency Injection)
	productRepo := repository.NewProductRepository(prismaClient)
	variantRepo := repository.NewVariantRepository(prismaClient)
	productService := service.NewProductService(productRepo, variantRepo, cacheClient)
	productHandler := handler.NewProductHandler(productService)

	// 5. Build HTTP router with route groups and rate limiters
	router := handler.SetupRouter(cfg, productHandler, cacheClient)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 6. Graceful shutdown handler
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("[INFO] Server listening on port %s (Environment: %s)", cfg.Port, cfg.Environment)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[FATAL] HTTP server error: %v", err)
		}
	}()

	<-shutdownChan
	log.Println("[INFO] Shutting down server gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("[ERROR] Server shutdown error: %v", err)
	}

	log.Println("[INFO] Server stopped gracefully.")
}
