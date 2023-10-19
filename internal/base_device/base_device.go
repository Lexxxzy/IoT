package base_device

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/Lexxxzy/iot-sockets/internal/encryption"
	"log"
)

var (
	mqttClient mqtt.Client
	clientID   string
)

func Initialize(deviceTag string) {
	opts := mqtt.NewClientOptions().AddBroker("tcp://mqttroute:1883").SetClientID(deviceTag)
	mqttClient = mqtt.NewClient(opts)
	clientID = deviceTag

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

func AcceptData(callback func(string)) {
	topic := clientID + "/#"
	mqttClient.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		decryptedData, err := encryption.Decrypt(msg.Payload())
		if err != nil {
			log.Fatal(err)
		}
		callback(string(decryptedData))
	})
}

func SendData(message []byte) {
	topic := clientID + "/commands"
	encryptedData := encryption.Encrypt(message)
	mqttClient.Publish(topic, 0, false, encryptedData)
}
