FROM golang:1.21 AS build-env

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/bot/main.go ./main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o telegram_bot .

FROM alpine
WORKDIR /app/
COPY --from=build-env /app/telegram_bot /app/telegram_bot
RUN apk add --no-cache bash

COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh
