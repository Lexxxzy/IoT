package main

import (
	"fmt"
	"github.com/Lexxxzy/iot-sockets/internal/base_device"
	"strings"
	"time"
)

const (
	deviceTag = "LIG"
	ON = "ON"
	OFF = "OFF"
	TOGGLE = "toggle"
)

var currentState string = "OFF"

func sendDataLightSwitch() {
	base_device.SendData([]byte(fmt.Sprintf("Light: %s", currentState)))
}

func handleLightSwitchData(data string) {
	data = strings.TrimRight(data, "\n")

	if data == TOGGLE {
		if currentState == ON {
			currentState = OFF
		} else {
			currentState = ON
		}
		sendDataLightSwitch()
		return
	}

	if (data == ON || data == OFF) && currentState != data {
		currentState = data
		sendDataLightSwitch()
	} else {
		fmt.Println("Invalid data!")
	}
}

func main() {
	base_device.Initialize(deviceTag)
	go base_device.AcceptData(handleLightSwitchData)

	for {
		sendDataLightSwitch()
		time.Sleep(30 * time.Second)
	}
}
