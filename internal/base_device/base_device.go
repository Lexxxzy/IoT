package base_device

import (
	"github.com/Lexxxzy/iot-sockets/internal/encryption"
	"log"
	"net"
)

func AcceptData(port string, callback func(string)) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(conn net.Conn) {
			defer conn.Close()

			data := make([]byte, 256)
			n, err := conn.Read(data)
			if err != nil {
				log.Println("Error reading from connection:", err)
				return
			}

			decryptedData, err := encryption.Decrypt(data[:n])
			if err != nil {
				log.Fatal(err)
			}
			callback(string(decryptedData))
		}(conn)
	}
}

func SendData(message []byte, deviceTag string) {
	conn, err := net.Dial("tcp", "gateway:8081")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	encryptedData := encryption.Encrypt(message)
	conn.Write(append([]byte("C"+deviceTag), encryptedData...))
}
