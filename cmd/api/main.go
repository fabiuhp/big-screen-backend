package main

import (
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

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		dbHost, dbPort, dbUser, dbPass, dbName)

	log.Printf("Tentando conectar ao banco de dados em: %s", dbHost)

	// Tenta conectar ao banco de dados com retry
	var db *sql.DB
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Tentativa %d: Erro ao conectar ao banco de dados: %v", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Testa a conexão
		err = db.Ping()
		if err != nil {
			log.Printf("Tentativa %d: Erro ao fazer ping no banco de dados: %v", i+1, err)
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
