package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/Lexxxzy/iot-sockets/internal/encryption"
	"log"
	"net"
    "strings"
)

var mqttClient mqtt.Client

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://mqttroute:1883").SetClientID("controller")
	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	mqttClient.Subscribe("#", 0, func(client mqtt.Client, msg mqtt.Message) {
		handleCommands(msg.Topic(), msg.Payload())
	})

	go startNetcatListener()

	select {}
}

func handleCommands(topic string, encryptedData []byte) {
	data, err := encryption.Decrypt(encryptedData)
	if err != nil {
		log.Printf("Decryption failed: %v", err)
		return
	}
	fmt.Printf("Received from %s: %s\n", topic, string(data))
}

func startNetcatListener() {
    ln, err := net.Listen("tcp", ":8085")
    if err != nil {
        log.Fatal(err)
    }
    defer ln.Close()

    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Println(err)
            continue
        }
        go handleNetcatConnection(conn)
    }
}

func handleNetcatConnection(conn net.Conn) {
    defer conn.Close()

    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        log.Println(err)
        return
    }

    inputData := string(buf[:n])
    parts := strings.SplitN(inputData, " ", 2)
    if len(parts) != 2 {
        log.Println("Invalid input format")
        return
    }

    deviceTag, command := parts[0], parts[1]
    topic := deviceTag + "/commands"
    sendCommands(topic, command)
}

func sendCommands(topic, command string) {
	encryptedCommand := encryption.Encrypt([]byte(command))
	mqttClient.Publish(topic, 0, false, encryptedCommand)
}
