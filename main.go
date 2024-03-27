package main

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

var botToken string
var commands = map[string]string{
	"/goodbye":      "Goodbye!",
	"/createwallet": "Create Wallet",
	"/listwallets":  "List Wallets",
}

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Read bot token from environment variable
	botToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
	}
}

func main() {
	port := ":8080"
	bot, err := tgbotapi.NewBotAPI(botToken)

	if err != nil {
		log.Panic(err)
	}

	updates := bot.ListenForWebhook("/webhook")

	go http.ListenAndServe(port, nil)

	log.Printf("Server is running and listening on port %s", port)

	processUpdates(updates, bot)
}

func processUpdates(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI) {
	for update := range updates {
		var msg tgbotapi.MessageConfig
		text := update.Message.Text

		log.Printf("Got message: %s", text)

		if update.Message == nil {
			continue
		}

		if text == "" {
			continue
		}

		switch text {
		case "/start":
			sendMenu(update.Message.Chat.ID, bot)
		default:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "I don't understand that command")
		}

		if commands[text] != "" {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, commands[text])
		}

		bot.Send(msg)
	}
}

func sendMenu(chatID int64, bot *tgbotapi.BotAPI) {
	// Create inline keyboard buttons
	btnCreateWallet := tgbotapi.NewInlineKeyboardButtonData("Create Wallet", "/createwallet")
	btnListWallets := tgbotapi.NewInlineKeyboardButtonData("List Wallets", "/listwallets")

	// Create a row for the buttons
	row := []tgbotapi.InlineKeyboardButton{btnCreateWallet, btnListWallets}
	rows := [][]tgbotapi.InlineKeyboardButton{row}

	// Create the inline keyboard markup
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	// Send message with inline keyboard to the user
	msg := tgbotapi.NewMessage(chatID, "Please choose an option:")
	msg.ReplyMarkup = inlineKeyboard

	bot.Send(msg)
}
