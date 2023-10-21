package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"github.com/eclipse/paho.mqtt.golang"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://mqttroute:1883")
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramBotToken == "" {
	  log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Fatal(err)
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatal(err)
	}

	sensorData := make(map[string]string)
	lastContact := time.Now()

	mqttClient.Subscribe("#", 0, func(client mqtt.Client, msg mqtt.Message) {
		lastContact = time.Now()
		switch msg.Topic() {
			case "HUM/commands":
				sensorData["Humidity"] = string(msg.Payload())
			case "THM/commands":
				sensorData["Thermostat"] = string(msg.Payload())
			case "LIG/commands":
				sensorData["Light"] = string(msg.Payload())
		}
	})

	for update := range updates {
		if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID
			callback := update.CallbackQuery.Data
			switch callback {
				case "trigger_light":
					mqttClient.Publish("LIG/commands", 0, false, "toggle")
					msg := tgbotapi.NewMessage(chatID, "Light switched!")
					bot.Send(msg)
				case "get_sensors":
					elapsed := time.Since(lastContact).Seconds()
					msgText := fmt.Sprintf("Last contact: %.0f seconds ago\n\n%s\n%s\n%s", elapsed, sensorData["Humidity"], sensorData["Thermostat"], sensorData["Light"])
					msg := tgbotapi.NewMessage(chatID, msgText)
					bot.Send(msg)

			}
			continue
		}

		inlineBtn1 := tgbotapi.NewInlineKeyboardButtonData("Trigger Light", "trigger_light")
		inlineBtn2 := tgbotapi.NewInlineKeyboardButtonData("Get Sensors Status", "get_sensors")
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(inlineBtn1),
			tgbotapi.NewInlineKeyboardRow(inlineBtn2),
		)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose an action:")
		msg.ReplyMarkup = inlineKeyboard
		bot.Send(msg)
	}
}
