# Используем официальный образ Go для сборки
FROM golang:1.24.1-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код проекта
COPY . .

# Компилируем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o marketplace ./cmd

# Создаем финальный образ
FROM alpine:latest

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем скомпилированный бинарник из builder
COPY --from=builder /app/marketplace .

# Копируем конфигурацию и миграции
COPY pkg/config/config.yaml ./pkg/config/config.yaml
COPY migrations ./migrations

# Указываем порт, который будет использовать приложение
EXPOSE 8080

# Команда для запуска приложения
CMD ["./marketplace"]