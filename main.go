package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/rafaeldajuda/sabor-e-arte-golang-telegram/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// mongodb
	// Configuração do cliente do MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://admin:admin@localhost:27017")

	// Conectando ao servidor do MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Verificando a conexão
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Conexão com MongoDB estabelecida com sucesso!")

	// Fechando a conexão com o banco de dados ao final do programa
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Conexão com MongoDB encerrada.")
	}()

	// Criando um documento BSON com os dados da imagem
	// photoBytes, err := os.ReadFile("./img/download.jpeg")
	// if err != nil {
	// 	panic(err)
	// }
	// menuItem := bson.M{
	// 	"nome":      "Pão de Queijo",
	// 	"tipo":      "Entrada",
	// 	"preco":     0.99,
	// 	"descricao": "Pão de queijo caseiro.",
	// 	"img":       photoBytes,
	// }
	// collection := client.Database("dev").Collection("sabor_arte")
	// result, err := collection.InsertOne(context.Background(), menuItem)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("menu_id", result)

	// pegando os items
	collection := client.Database("dev").Collection("sabor_arte")
	result, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal("find", err)
	}

	for result.Next(context.Background()) {
		raw := result.Current
		itemMenu := entity.ItemMenu{}
		err := bson.UnmarshalExtJSON([]byte(raw.String()), false, &itemMenu)
		if err != nil {
			log.Fatal("LOOP", err)
		}

		fmt.Println("item", itemMenu)
		fmt.Println("img", itemMenu.Imagem)
	}

	// bot telegram
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
