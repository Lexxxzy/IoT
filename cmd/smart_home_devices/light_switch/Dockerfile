# Build stage
FROM golang:1.21 AS build-env
WORKDIR /app
COPY go.mod .
COPY go.sum . 
COPY vendor/ ./vendor/
COPY internal/ ./internal/
COPY cmd/smart_home_devices/light_switch/light_switch.go ./light_switch.go
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o light_switch .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=build-env /app/light_switch /app/light_switch
CMD ["/app/light_switch"]