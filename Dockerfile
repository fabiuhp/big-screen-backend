FROM golang:1.21-alpine AS builder

WORKDIR /app

# Instala as dependências necessárias
RUN apk add --no-cache gcc musl-dev

# Copia os arquivos de dependências
COPY go.mod ./
COPY go.sum ./

# Baixa as dependências
RUN go mod download

# Copia o código fonte
COPY . .

# Compila a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Imagem final
FROM alpine:latest

WORKDIR /app

# Copia o binário compilado
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Expõe a porta da aplicação
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./main"] 