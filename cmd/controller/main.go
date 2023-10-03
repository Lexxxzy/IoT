package main

import (
	"fmt"
	"github.com/Lexxxzy/iot-sockets/internal/encryption"
	"log"
	"net"
)

func handleCommands(encryptedData []byte) {
	data, err := encryption.Decrypt(encryptedData)
	if err != nil {
		log.Printf("Decryption failed: %v", err)
		return
	}
	fmt.Println("Received:", string(data))
}

func sendCommands(command string) {
	encryptedCommand := encryption.Encrypt([]byte(command[4:]))

	gtw, err := net.Dial("tcp", "gateway:8081")
	if err != nil {
		log.Printf("Error dialing controller: %v", err)
	}

	_, err = gtw.Write(append([]byte(command[:4]), encryptedCommand...))
	if err != nil {
		log.Printf("Could not write to gateway: %v", err)
	}
}

func main() {
	// Listener for incoming data
	go func() {
		listener, err := net.Listen("tcp", ":8080")
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
				handleCommands(buf[:n])
			}(conn)
		}
	}()

	// Listener for outgoing data
	go func() {
		outgoingListener, err := net.Listen("tcp", ":8085")
		if err != nil {
			log.Fatalf("Error setting up outgoing listener: %v", err)
		}
		for {
			conn, err := outgoingListener.Accept()
			if err != nil {
				log.Printf("Error accepting connection on outgoing port: %v", err)
				continue
			}
			go func(conn net.Conn) {
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					log.Printf("Error reading from connection: %v", err)
					return
				}
				command := string(buf[:n])
				sendCommands(command)
			}(conn)
		}
	}()

	select {}
}
