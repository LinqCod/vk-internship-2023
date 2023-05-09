package telegrambot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"log"
)

func initKeyboards() []tgbotapi.ReplyKeyboardMarkup {
	keyboards := make([]tgbotapi.ReplyKeyboardMarkup, 5)

	keyboards[0] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("math"),
			tgbotapi.NewKeyboardButton("trivia"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("date"),
			tgbotapi.NewKeyboardButton("year"),
		),
	)
	keyboards[1] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("About specific number"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("About random number"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("<- Back"),
		),
	)
	keyboards[2] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("About specific number"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("About random number"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("<- Back"),
		),
	)
	keyboards[3] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("About specific date"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("About random date"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("<- Back"),
		),
	)
	keyboards[4] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("About specific year"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("About random year"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("<- Back"),
		),
	)

	return keyboards
}

type NumbersBot struct {
	bot          *tgbotapi.BotAPI
	updateConfig tgbotapi.UpdateConfig
	keyboards    []tgbotapi.ReplyKeyboardMarkup
}

func NewNumbersBot() (*NumbersBot, error) {
	token := viper.GetString("TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	offset := viper.GetInt("UPDATE_OFFSET")
	timeout := viper.GetInt("UPDATE_TIMEOUT")
	u := tgbotapi.UpdateConfig{
		Offset:  offset,
		Timeout: timeout,
	}

	keyboards := initKeyboards()

	return &NumbersBot{
		bot:          bot,
		updateConfig: u,
		keyboards:    keyboards,
	}, nil
}

func (b *NumbersBot) Start() {
	updates, err := b.bot.GetUpdatesChan(b.updateConfig)
	if err != nil {
		log.Fatal(err)
	}

	currentCommandSection := ""
	for update := range updates {
		if update.Message == nil {
			continue
		}

		msgText := ""
		keyboardId := 0
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		switch update.Message.Text {
		case "/start":
			msgText = "Hello, glad to see u with us :)\nLet's explore some interesting facts about numbers!"
			currentCommandSection = ""
		case "math":
			msgText = "Choose one option to see a math fact!"
			currentCommandSection = "math"
			keyboardId = 1
		case "trivia":
			msgText = "Choose one option to see a trivia fact!"
			currentCommandSection = "trivia"
			keyboardId = 2
		case "date":
			msgText = "Choose one option to see a specific date fact!"
			currentCommandSection = "date"
			keyboardId = 3
		case "year":
			msgText = "Choose one option to see a specific year fact!"
			currentCommandSection = "year"
			keyboardId = 4
		case "<- Back":
			msgText = "what would you like to choose then?"
			currentCommandSection = ""
		default:
			msgText = "I don't know this command :("
			currentCommandSection = ""
		}

		log.Println(currentCommandSection)

		msg.ReplyMarkup = b.keyboards[keyboardId]
		msg.Text = msgText
		if _, err := b.bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
