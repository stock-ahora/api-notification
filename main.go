// cmd/api/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"api-notification/internal/handlers"
)

func main() {
	fmt.Println("Servidor iniciando...")

	// Configurar rutas
	setupRoutes()

	// Crear servidor HTTP
	server := &http.Server{
		Addr: ":8084",
	}

	// Iniciar servidor
	fmt.Println("Escuchando en el puerto 8084")
	log.Fatal(server.ListenAndServe())
}

func setupRoutes() {
	http.HandleFunc("/health", handlers.HealthCheck)
}
