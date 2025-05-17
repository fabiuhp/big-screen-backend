package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	// Log das variáveis de ambiente (sem a senha)
	log.Printf("Configuração do banco de dados:")
	log.Printf("Host: %s", dbHost)
	log.Printf("Port: %s", dbPort)
	log.Printf("User: %s", dbUser)
	log.Printf("Database: %s", dbName)

	// String de conexão com SSL requerido
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require",
		dbUser, dbPass, dbHost, dbPort, dbName)

	log.Printf("Tentando conectar ao banco de dados...")

	// Tenta conectar ao banco de dados com retry
	var db *sql.DB
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		log.Printf("Tentativa %d de %d", i+1, maxRetries)

		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Erro ao abrir conexão: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Testa a conexão
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = db.PingContext(ctx)
		if err != nil {
			log.Printf("Erro ao fazer ping: %v", err)
			db.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("Conexão com o banco de dados estabelecida com sucesso!")
		break
	}

	if err != nil {
		log.Fatal("Não foi possível conectar ao banco de dados após várias tentativas")
	}
	defer db.Close()

	// Configura o pool de conexões
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

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
