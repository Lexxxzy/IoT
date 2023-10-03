# Build stage
FROM golang:1.21 AS build-env
WORKDIR /app
COPY go.mod .
COPY internal/ ./internal/
COPY cmd/smart_home_devices/thermostat/thermostat.go ./thermostat.go
RUN CGO_ENABLED=0 GOOS=linux go build -o thermostat .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=build-env /app/thermostat /app/thermostat
CMD ["/app/thermostat"]