package main

import (
	"encoding/json"
	"fmt"
	"github.com/Lexxxzy/iot-sockets/internal/base_device"
	"github.com/labstack/echo/v4"
	"time"
)

const deviceTag = "THM"

func sendDataThermostat() {
	var data = map[string]string{"device": "Thermostat", "state": fmt.Sprintf("%d°C", 20)}
	newData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	base_device.SendData(newData)
}

func handleThermostatData(data string) {
	fmt.Println("Received data for thermostat:", data)
}

func main() {
	e := echo.New()

	base_device.Initialize(deviceTag)
	base_device.AcceptData(e, handleThermostatData)

	go func() {
		for {
			sendDataThermostat()
			time.Sleep(30 * time.Second)
		}
	}()

	e.Start(":8080")
}
