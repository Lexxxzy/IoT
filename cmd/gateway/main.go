package main

import (
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	dataChannel := make(chan []byte)

	go func() {
		listener, err := net.Listen("tcp", ":8081")
		if err != nil {
			log.Fatalf("Error setting up listener: %v", err)
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Error accepting connection: %v", err)
				continue
			}
			go func(conn net.Conn) {
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					log.Printf("Error reading from connection: %v", err)
					return
				}
				dataChannel <- buf[:n]
			}(conn)
		}
	}()

	for {
		aggregatedData := <-dataChannel
		direction := string(aggregatedData[0])    // Первый байт - направление
		deviceType := string(aggregatedData[1:4]) // Другие 3 байта - тип устройства
		payload := aggregatedData[4:]

		target := ""
		switch strings.ToUpper(deviceType) {
		case "THM":
			target = "thermostat:8085"
		case "HUM":
			target = "humidity_sensor:8085"
		case "LIG":
			target = "light_switch:8085"
		}

		if strings.ToUpper(direction) == "C" {
			// Send to controller
			conn, err := net.Dial("tcp", "controller:8080")
			if err != nil {
				log.Printf("Error dialing controller: %v", err)
				continue
			}
			conn.Write(payload)
			conn.Close()
		} else if strings.ToUpper(direction) == "D" && target != "" {
			conn, err := net.Dial("tcp", target)
			if err != nil {
				log.Printf("Error dialing device: %v", err)
				continue
			}
			conn.Write(payload)
			conn.Close()
		}

		time.Sleep(1 * time.Second)
	}
}
