package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var menu = map[string]string{
	"1": "Menu",
	"2": "Ver mesas disponíveis",
	"3": "Fazer reserva",
	"4": "Minhas reservas",
	"5": "Cancelar reserva",
}

var menuKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1. Menu"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("2. Ver mesas disponíveis"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("3. Fazer reserva"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("4. Minhas reservas"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("5. Cancelar reserva"),
	),
)

func main() {
	godotenv.Load()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := bot.GetUpdatesChan(updateConfig)

	// for update := range updates {
	// 	if update.Message == nil {
	// 		continue
	// 	}

	// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	// 	msg.ReplyToMessageID = update.Message.MessageID

	// 	if _, err := bot.Send(msg); err != nil {
	// 		panic(err)
	// 	}
	// }

	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		t := tgbotapi.BotCommand{Command: "start", Description: "Seja bem vindo ao restaurante Sabor e Arte. Para ver as opções digite /open"}

		switch update.Message.Command() {
		case "/start":
			msg.Text = t.Description
		case "/open":
			msg.ReplyMarkup = menuKeyboard
		case "/close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		default:
			msg.Text = "Olá, tudo bem? Para ver as opções digite /start"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}

}
