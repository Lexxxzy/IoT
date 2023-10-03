package main

import (
	"fmt"
	"github.com/Lexxxzy/iot-sockets/internal/base_device"
	"time"
)

const deviceTag = "THM"

func sendDataThermostat() {
	base_device.SendData([]byte(fmt.Sprintf("Thermostat: %d°C", 20)), deviceTag)
}

func handleThermostatData(data string) {
	fmt.Println("Received data for thermostat:", data)
}

func main() {
	go base_device.AcceptData(":8085", handleThermostatData)

	for {
		sendDataThermostat()
		time.Sleep(30 * time.Second)
	}
}
