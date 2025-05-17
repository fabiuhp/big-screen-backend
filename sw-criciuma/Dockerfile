FROM golang:1.21-alpine

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
RUN go build -o main ./cmd/api

# Expõe a porta da aplicação
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./main"] 