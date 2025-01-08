# Используем официальный образ Go
FROM golang:1.22.4 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go mod и sum файлы
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o tabhub ./cmd/tabhub

# Создаем финальный легковесный образ
FROM alpine:latest

WORKDIR /root/

# Копируем собранный бинарный файл из builder образа
COPY --from=builder /app/tabhub .

# Копируем .env файл (если есть)
COPY .env* ./

# Устанавливаем порт, который будет использоваться приложением
EXPOSE ${APP_PORT}

# Команда для запуска приложения
ENTRYPOINT ["./tabhub"]