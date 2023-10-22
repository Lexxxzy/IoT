package main

import (
	"encoding/json"
	"fmt"
	"github.com/Lexxxzy/iot-sockets/internal/base_device"
	"github.com/labstack/echo/v4"
	"time"
)

const deviceTag = "HUM"

func sendDataHumiditySensor() {
	var data = map[string]string{"device": "Humidity", "state": fmt.Sprintf("%d%%", 60)}
	newData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	base_device.SendData(newData)
}

func handleHumiditySensorData(data string) {
	fmt.Println("Received data for humidity sensor:", data)
}

func main() {
	e := echo.New()

	base_device.Initialize(deviceTag)
	base_device.AcceptData(e, handleHumiditySensorData)

	go func() {
		for {
			sendDataHumiditySensor()
			time.Sleep(30 * time.Second)
		}
	}()

	e.Start(":8080")
}
