package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Lexxxzy/iot-sockets/internal/encryption"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"sync"
)

var deviceURLs = map[string]string{
	"LIG": "http://light_switch:8080/command",
	"THM": "http://thermostat:8080/command",
	"HUM": "http://humidity_sensor:8080/command",
}

var mutex = &sync.RWMutex{}

type Container struct {
	mu   sync.Mutex
	data map[string]string
}

func (c *Container) renewStatus(jsonData map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[jsonData["device"]] = jsonData["state"]
}

var dataContainer = Container{
	data: map[string]string{
		"Light":      "N/A",
		"Thermostat": "N/A",
		"Humidity":   "N/A",
	},
}

func main() {
	e := echo.New()

	e.POST("/receive", handleReceive)
	e.POST("/send/:deviceID", handleSend)
	e.GET("/", getData)

	e.Start(":8080")
}

func handleReceive(c echo.Context) error {
	encryptedData, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
	}

	received, err := encryption.Decrypt(encryptedData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Decryption failed"})
	}

	var jsonData = map[string]string{}
	err = json.Unmarshal(received, &jsonData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	dataContainer.renewStatus(jsonData)
	fmt.Printf("Received: %s\n", string(received))

	return c.JSON(http.StatusOK, dataContainer.data)
}

func getData(c echo.Context) error {
	return c.JSON(http.StatusOK, dataContainer.data)
}

func handleSend(c echo.Context) error {
	deviceID := c.Param("deviceID")
	commandData, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
	}
	println("DEBUG: ", string(commandData))

	deviceURL, exists := deviceURLs[deviceID]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Device not found"})
	}

	encryptedCommand := encryption.Encrypt(commandData)
	response, err := http.Post(deviceURL, "application/octet-stream", bytes.NewReader(encryptedCommand))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send command"})
	}
	println("DEBUG: ", response.Status)
	defer response.Body.Close()

	return c.JSON(http.StatusOK, map[string]string{"status": "sent"})
}
