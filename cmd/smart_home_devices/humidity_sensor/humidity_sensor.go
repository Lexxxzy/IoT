package main

import (
	"fmt"
	"github.com/Lexxxzy/iot-sockets/internal/base_device"
	"time"
)

const deviceTag = "HUM"

func sendDataHumiditySensor() {
	base_device.SendData([]byte(fmt.Sprintf("Humidity: %d%%", 60)), deviceTag)
}

func handleHumiditySensorData(data string) {
	fmt.Println("Received data for humidity sensor:", data)
}

func main() {
	go base_device.AcceptData(":8085", handleHumiditySensorData)

	for {
		sendDataHumiditySensor()
		time.Sleep(30 * time.Second)
	}
}
