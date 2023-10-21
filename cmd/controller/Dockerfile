# Build stage
FROM golang:1.21 AS build-env
WORKDIR /app
COPY go.mod .
COPY go.sum . 
#COPY vendor/ ./vendor/
COPY internal/ ./internal/
COPY cmd/controller/main.go ./main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o controller .

# Final stage
FROM ubuntu
WORKDIR /app
COPY --from=build-env /app/controller /app/controller
COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh