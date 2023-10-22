# Build stage
FROM golang:1.21 AS build-env
WORKDIR /app
COPY go.mod .
COPY go.sum .
COPY internal/ ./internal/
COPY cmd/smart_home_devices/humidity_sensor/humidity_sensor.go ./humidity_sensor.go
RUN CGO_ENABLED=0 GOOS=linux go build -o humidity_sensor .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=build-env /app/humidity_sensor /app/humidity_sensor
CMD ["/app/humidity_sensor"]