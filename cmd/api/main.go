package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	httpHandler "github.com/abiopereira/sw-criciuma/internal/delivery/http"
	"github.com/abiopereira/sw-criciuma/internal/repository/postgres"
	"github.com/abiopereira/sw-criciuma/internal/usecase"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Configuração do banco de dados
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASS", "postgres")
	dbName := getEnv("DB_NAME", "sw_criciuma")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Inicialização das dependências
	messageRepo := postgres.NewMessageRepository(db)
	messageUseCase := usecase.NewMessageUseCase(messageRepo)

	// Configuração do router
	router := mux.NewRouter()
	httpHandler.NewMessageHandler(router, messageUseCase)

	// Configuração do servidor
	port := getEnv("PORT", "8080")
	log.Printf("Servidor iniciado na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
