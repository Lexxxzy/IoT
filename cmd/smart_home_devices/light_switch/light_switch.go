package main

import (
	"encoding/json"
	"fmt"
	"github.com/Lexxxzy/iot-sockets/internal/base_device"
	"github.com/labstack/echo/v4"
	"strings"
	"time"
)

const deviceTag = "LIG"

var currentState string = "OFF"

func sendDataLightSwitch() {
	var data = map[string]string{"device": "Light", "state": currentState}
	newData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	base_device.SendData(newData)
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
	e := echo.New()

	base_device.Initialize(deviceTag)
	base_device.AcceptData(e, handleLightSwitchData)

	go func() {
		for {
			sendDataLightSwitch()
			time.Sleep(30 * time.Second)
		}
	}()

	e.Start(":8080")
}
