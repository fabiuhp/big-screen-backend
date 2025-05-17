# SW Criciúma - Backend

Backend do sistema de mensagens para o telão do SW Criciúma, desenvolvido em Go usando Clean Architecture.

## Requisitos

- Go 1.21 ou superior
- PostgreSQL 12 ou superior
- Docker e Docker Compose (opcional)

## Configuração

### Executando Localmente

1. Clone o repositório
2. Instale as dependências:
```bash
go mod download
```

3. Configure as variáveis de ambiente (ou use os valores padrão):
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASS=postgres
export DB_NAME=sw_criciuma
export PORT=8080
```

4. Crie o banco de dados e execute as migrações:
```bash
psql -U postgres -c "CREATE DATABASE sw_criciuma;"
psql -U postgres -d sw_criciuma -f migrations/001_create_messages_table.sql
```

5. Execute o projeto:
```bash
go run cmd/api/main.go
```

### Executando com Docker

1. Clone o repositório
2. Execute o comando:
```bash
docker-compose up --build
```

A aplicação estará disponível em `http://localhost:8080`

## Endpoints

### GET /api/messages
Retorna todas as mensagens na ordem em que foram enviadas (FIFO).

### POST /api/messages
Envia uma nova mensagem.

Body:
```json
{
  "type": "text",
  "content": "Texto da mensagem",
  "duration": 10
}
```

## Estrutura do Projeto

```
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/
│   │   └── message.go
│   ├── usecase/
│   │   └── message_usecase.go
│   ├── repository/
│   │   └── postgres/
│   │       └── message_repository.go
│   └── delivery/
│       └── http/
│           └── message_handler.go
├── migrations/
│   └── 001_create_messages_table.sql
├── Dockerfile
├── docker-compose.yml
└── go.mod
``` 