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

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/stock-ahora/api-notification/internal/config"
	"github.com/stock-ahora/api-notification/internal/messaging"
)

func main() {
	// Cargar variables de entorno desde el archivo .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	// 1. Cargar configuración
	cfg, err := config.LoadSecrets()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	log.Printf("Starting API on port %d", cfg.ServerPort)

	// 2. Crear consumer de RabbitMQ
	consumer, err := messaging.NewConsumer(*cfg.GetRabbitMQConfig())
	if err != nil {
		log.Fatal("Error creating consumer:", err)
	}
	defer consumer.Close()

	// 3. Crear router
	router := mux.NewRouter()

	// 4. Configurar rutas
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// 5. Configurar servidor HTTP
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	// 6. Iniciar servidor
	go func() {
		log.Printf("Server listening on :%d", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error:", err)
		}
	}()

	// 7. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
