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
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Ex6linz/BookSwap/backend/internal/auth/repository/postgres"
	authService "github.com/Ex6linz/BookSwap/backend/internal/auth/service"
	authRest "github.com/Ex6linz/BookSwap/backend/internal/auth/transport/rest"
	"github.com/Ex6linz/BookSwap/backend/pkg/config"
)

func main() {
	// 1. Ładowanie konfiguracji
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("Błąd ładowania konfiguracji: %v", err)
	}

	// 2. Inicjalizacja połączenia z bazą danych
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Nie można połączyć się z bazą danych: %v", err)
	}
	defer dbPool.Close()

	// 3. Inicjalizacja komponentów autentykacji
	userRepo := postgres.NewUserRepository(dbPool)
	authSvc := authService.NewAuthService(userRepo, cfg.JWT.Secret)
	authHandler := authRest.NewAuthHandler(authSvc)

	// 4. Konfiguracja routera Gin
	router := gin.Default()

	// 5. Rejestracja endpointów
	// Publiczne endpointy
	public := router.Group("/api/v1")
	{
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
	}

	// Chronione endpointy (wymagają JWT)
	protected := router.Group("/api/v1")
	protected.Use(authRest.AuthMiddleware(cfg.JWT.Secret))
	{
		// Tutaj dodamy później chronione endpointy
		protected.GET("/test-auth", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Jesteś zalogowany!"})
		})
	}

	// 6. Konfiguracja serwera HTTP
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// 7. Uruchomienie serwera w goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Błąd serwera: %v", err)
		}
	}()

	log.Printf("Serwer działa na porcie %s", cfg.Server.Port)

	// 8. Obsługa sygnałów do poprawnego zamknięcia
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Zamykanie serwera...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Wymuszone zamknięcie serwera:", err)
	}

	log.Println("Serwer zatrzymany")
}
