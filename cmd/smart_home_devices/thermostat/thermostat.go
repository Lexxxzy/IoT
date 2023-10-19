package main

import (
	"fmt"
	"github.com/Lexxxzy/iot-sockets/internal/base_device"
	"time"
)

const deviceTag = "THM"

func sendDataThermostat() {
	base_device.SendData([]byte(fmt.Sprintf("Thermostat: %d°C", 20)))
}

func handleThermostatData(data string) {
	fmt.Println("Received data for thermostat:", data)
}

func main() {
	base_device.Initialize(deviceTag)
	go base_device.AcceptData(handleThermostatData)

	for {
		sendDataThermostat()
		time.Sleep(30 * time.Second)
	}
}
