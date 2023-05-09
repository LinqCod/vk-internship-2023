package telegrambot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/linqcod/vk-internship-2023/internal/numbersapi"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"strings"
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
	numbersApi   *numbersapi.Api
}

func NewNumbersBot(numbersApi *numbersapi.Api) (*NumbersBot, error) {
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
		numbersApi:   numbersApi,
	}, nil
}

func (b *NumbersBot) Start() {
	updates, err := b.bot.GetUpdatesChan(b.updateConfig)
	if err != nil {
		log.Fatal(err)
	}

	currentCategory := ""
	prevCommand := ""
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
			prevCommand = ""
			currentCategory = ""
		case "math":
			msgText = "Choose one option to see a math fact!"
			keyboardId = 1
			currentCategory = "math"
			prevCommand = ""
		case "trivia":
			msgText = "Choose one option to see a trivia fact!"
			keyboardId = 2
			currentCategory = "trivia"
			prevCommand = ""
		case "About specific number":
			if currentCategory == "math" || currentCategory == "trivia" {
				msgText = "Please, enter the number to get fact about it:"
				prevCommand = "About specific number"
			} else {
				msgText = "Please, choose the category!"
			}
		case "About random number", "About random date", "About random year":
			if prevCommand == "About specific number" {
				msgText = "First, u need to complete your specific number fact request! So enter your number:"
			} else if prevCommand == "About specific date" {
				msgText = "First, u need to complete your specific date fact request! So enter your date in format (month/day):"
			} else if prevCommand == "About specific year" {
				msgText = "First, u need to complete your specific year fact request! So enter your year:"
			} else if currentCategory == "math" || currentCategory == "trivia" || currentCategory == "date" || currentCategory == "year" {
				msgText, err = b.numbersApi.GetFact("random", currentCategory)
				if err != nil {
					log.Fatalf("error while getting fact: %s", err)
				}
				currentCategory = ""
				prevCommand = ""
			} else {
				msgText = "Please, choose the category!"
			}
		case "About specific date":
			if currentCategory == "date" {
				msgText = "Please, enter the date in format (month/day without leading zeros) to get fact about it:"
				prevCommand = "About specific date"
			} else {
				msgText = "Please, choose the category!"
			}
		case "About specific year":
			if currentCategory == "year" {
				msgText = "Please, enter the year to get fact about it:"
				prevCommand = "About specific year"
			} else {
				msgText = "Please, choose the category!"
			}
		case "date":
			msgText = "Choose one option to see a specific date fact!"
			keyboardId = 3
			currentCategory = "date"
			prevCommand = ""
		case "year":
			msgText = "Choose one option to see a specific year fact!"
			keyboardId = 4
			currentCategory = "year"
			prevCommand = ""
		case "<- Back":
			msgText = "what would you like to choose then?"
			currentCategory = ""
			prevCommand = ""
		default:
			if prevCommand == "About specific year" || prevCommand == "About specific number" {
				if _, err = strconv.ParseInt(update.Message.Text, 10, 32); err != nil {
					log.Println(err)
					msgText = "Please, enter a valid number!"
				} else {
					msgText, err = b.numbersApi.GetFact(update.Message.Text, currentCategory)
					if err != nil {
						log.Fatalf("error while getting fact: %s", err)
					}

					currentCategory = ""
					prevCommand = ""
				}
			} else if prevCommand == "About specific date" {
				monthAndDay := strings.Split(update.Message.Text, "/")
				if len(monthAndDay) == 2 && monthAndDay[0] != "" && monthAndDay[1] != "" {
					if _, err = strconv.ParseInt(monthAndDay[0], 10, 32); err != nil {
						log.Println(err)
						msgText = "Please, enter a valid month!"
					} else if _, err = strconv.ParseInt(monthAndDay[1], 10, 32); err != nil {
						log.Println(err)
						msgText = "Please, enter a valid day!"
					} else {
						msgText, err = b.numbersApi.GetFact(update.Message.Text, currentCategory)
						if err != nil {
							log.Fatalf("error while getting fact: %s", err)
						}

						currentCategory = ""
						prevCommand = ""
					}
				} else {
					msgText = "Please, enter a valid date!"
				}
			} else {
				msgText = "I don't know this command :("
				currentCategory = ""
				prevCommand = ""
			}
		}

		msg.ReplyMarkup = b.keyboards[keyboardId]
		msg.Text = msgText
		if _, err := b.bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
