package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/stock-ahora/api-notification/internal/config"
	"github.com/stock-ahora/api-notification/internal/handlers"
	"github.com/stock-ahora/api-notification/internal/messaging"
	"github.com/stock-ahora/api-notification/internal/service"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Conectar a RabbitMQ
	mqConn, mqChannel, err := messaging.Connect(cfg.RabbitMQ)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer mqConn.Close()
	defer mqChannel.Close()

	// Crear publisher
	publisher := messaging.NewPublisher(mqChannel)

	// Crear servicio
	notificationService := service.NewNotificationService(publisher)

	// Iniciar consumer
	consumer := messaging.NewConsumer(mqChannel, notificationService)
	go consumer.Start()

	// Configurar HTTP
	httpHandlers := handlers.NewHandlers(notificationService)
	router := handlers.SetupRouter(httpHandlers)

	// Servidor HTTP
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Graceful shutdown
	go gracefulShutdown(srv, consumer)

	log.Printf("🚀 Server starting on port %s", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to start server:", err)
	}
}

func gracefulShutdown(srv *http.Server, consumer *messaging.Consumer) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("🛑 Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	consumer.Stop()
	srv.Shutdown(ctx)
}
