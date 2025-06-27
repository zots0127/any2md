package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"any2md/internal/adapters/handlers"
	"any2md/internal/infrastructure/config"
	"any2md/internal/infrastructure/middleware"
	"any2md/internal/usecases"
)

func main() {
	cfg := config.Load()
	
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimiter(cfg.RateLimit.MaxRequests, cfg.RateLimit.Window))
	
	converterUseCase := usecases.NewConverterUseCase()
	httpHandler := handlers.NewHTTPHandler(converterUseCase)
	
	router.GET("/health", httpHandler.Health)
	router.POST("/api/v1/convert", httpHandler.Convert)
	
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	
	go func() {
		fmt.Printf("Server starting on port %s...\n", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	fmt.Println("Shutting down server...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	
	fmt.Println("Server exited")
}