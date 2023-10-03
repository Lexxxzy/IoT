# Build stage
FROM golang:1.21 AS build-env
WORKDIR /app
COPY go.mod .
COPY internal/ ./internal/
COPY cmd/gateway/main.go ./main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway .

# Final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /app/gateway /app/gateway
CMD ["/app/gateway"]
