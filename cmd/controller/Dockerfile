# Build stage
FROM golang:1.21 AS build-env
WORKDIR /app
COPY go.mod .
COPY internal/ ./internal/
COPY cmd/controller/main.go ./main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o controller .

# Final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /app/controller /app/controller
CMD ["/app/controller"]
