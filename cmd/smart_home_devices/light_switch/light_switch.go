package main

import (
	"fmt"
	"github.com/Lexxxzy/iot-sockets/internal/base_device"
	"strings"
	"time"
)

const deviceTag = "LIG"

var currentState string = "OFF"

func sendDataLightSwitch() {
	base_device.SendData([]byte(fmt.Sprintf("Light: %s", currentState)), deviceTag)
}

func handleLightSwitchData(data string) {
	data = strings.TrimRight(data, "\n")
	if (data == "ON" || data == "OFF") && currentState != data {
		currentState = data
		sendDataLightSwitch()
	} else {
		fmt.Println("Invalid data!")
	}
}

func main() {
	go base_device.AcceptData(":8085", handleLightSwitchData)

	for {
		sendDataLightSwitch()
		time.Sleep(30 * time.Second)
	}
}
