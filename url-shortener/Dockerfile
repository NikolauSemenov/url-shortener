FROM golang:1.24-alpine

WORKDIR /app

# Устанавливаем зависимости + инструменты для разработки
RUN apk add --no-cache git gcc musl-dev bash

# Копируем файлы модулей
COPY go.mod go.sum ./
COPY db ./
RUN go mod download
