package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/rafaeldajuda/sabor-e-arte-golang-telegram/entity"
)

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

var menuRestaurante entity.MenuRestaurante

func init() {
	item1 := entity.ItemMenu{
		Nome:      "Pão de Queijo",
		Preco:     0.99,
		Tipo:      "Entrada",
		Descricao: "Pão de queijo caseiro.",
		Imagem:    []byte{},
	}
	item2 := entity.ItemMenu{
		Nome:      "Arroz",
		Preco:     5.99,
		Tipo:      "Comida",
		Descricao: "Arroz caseiro.",
		Imagem:    []byte{},
	}
	item3 := entity.ItemMenu{
		Nome:      "Sorvete",
		Preco:     2.99,
		Tipo:      "Sobremesa",
		Descricao: "Sorvete caseiro.",
		Imagem:    []byte{},
	}
	menuRestaurante = append(menuRestaurante, item1, item2, item3)
}

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
		// t := tgbotapi.BotCommand{Command: "start", Description: "Seja bem vindo ao restaurante Sabor e Arte. Para ver as opções digite /open \n\nOBS: para fechar as opções digite /close"}

		option := func(cmd string, msg string) string {
			if cmd == "" {
				return msg
			}
			return cmd
		}

		switch option(update.Message.Command(), update.Message.Text) {
		case "start":
			msg.Text = "Seja bem vindo ao restaurante Sabor e Arte. Para ver as opções digite /open \n\nOBS: para fechar as opções digite /close"
		case "open":
			msg.ReplyMarkup = menuKeyboard
		case "close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case "1. Menu":
			for _, v := range menuRestaurante {
				text := fmt.Sprintf("Nome: %s\nPreço: %.2f\nTipo: %s\n\nDescrição: %s", v.Nome, v.Preco, v.Tipo, v.Descricao)
				img := "./img/download.jpeg"
				photoBytes, err := os.ReadFile(img)
				if err != nil {
					panic(err)
				}
				photoTele := tgbotapi.FileBytes{
					Name:  "comida",
					Bytes: photoBytes,
				}

				msg := tgbotapi.NewPhoto(update.Message.Chat.ID, photoTele)
				msg.Caption = text
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}
			continue
		case "2. Ver mesas disponíveis":
			msg.Text = "Ver mesas disponíveis"
		case "3. Fazer reserva":
			msg.Text = "Fazer reserva"
		case "4. Minhas reservas":
			msg.Text = "Minhas reservas"
		case "5. Cancelar reserva":
			msg.Text = "Cancelar reserva"
		default:
			msg.Text = "Olá, tudo bem? Para ver as opções digite /start"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}

}
