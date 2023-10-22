package base_device

import (
	"bytes"
	"github.com/Lexxxzy/iot-sockets/internal/encryption"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	clientID string
)

func Initialize(deviceTag string) {
	clientID = deviceTag
}

func AcceptData(e *echo.Echo, callback func(string)) {
	e.POST("/command", func(c echo.Context) error {
		encryptedData, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			log.Fatal(err)
		}
		decryptedData, err := encryption.Decrypt(encryptedData)
		if err != nil {
			log.Fatal(err)
		}
		callback(string(decryptedData))
		return c.JSON(http.StatusOK, map[string]string{"status": "success"})
	})
}

func SendData(message []byte) int {
	endpoint := "http://controller:8080/receive"
	encryptedData := encryption.Encrypt(message)
	_, err := http.Post(endpoint, "application/octet-stream", bytes.NewReader(encryptedData))
	if err != nil {
		log.Fatal(err)
	}

	return 200
}
