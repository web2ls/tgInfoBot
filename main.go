package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type exactTime struct {
	Time int `json:"time"`
}

type currencyItem struct {
	ID       string
	NumCode  string
	CharCode string
	Nominal  int
	Name     string
	Value    float32
	Previous float32
}

type exhanges struct {
	Valute map[string]currencyItem
}

var buttonsMarkup = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Точное время (московское)"),
		tgbotapi.NewKeyboardButton("Курсы валют"),
	),
)

const timeURL = "https://yandex.com/time/sync.json?geo=213"
const currencyURL = "https://www.cbr-xml-daily.ru/daily_json.js"

const msInSec = 1000

func checkOnAdditionalZero(dateValue string) string {
	if len(dateValue) < 2 {
		return "0" + dateValue
	}

	return dateValue
}

func getTimeAsString(timestamp int) string {
	fullDate := time.Unix(int64(timestamp/msInSec), 0)
	location, _ := time.LoadLocation("Europe/Moscow")
	fullDateByLocation := fullDate.In(location)
	hour := checkOnAdditionalZero(strconv.Itoa(fullDateByLocation.Hour()))
	minute := checkOnAdditionalZero(strconv.Itoa(fullDateByLocation.Minute()))
	second := checkOnAdditionalZero(strconv.Itoa(fullDateByLocation.Second()))

	return fmt.Sprintf("%s : %s : %s", string(hour), string(minute), string(second))
}

func getCurrencyValue(data exhanges, mainCurrency []string) string {
	var usdValue float32
	var eurValue float32
	for _, value := range mainCurrency {
		if value == "USD" {
			usdValue = data.Valute[value].Value
		} else if value == "EUR" {
			eurValue = data.Valute[value].Value
		}
		fmt.Println(value)
	}

	return fmt.Sprintf("$%.2f ¢%.2f", usdValue, eurValue)
}

func main() {
	fmt.Println("Bot has been started...")

	TOKEN := os.Getenv("TOKEN")
	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg.Text = "You are wellcome"

			default:
				msg.Text = "Unknown command, sorry"
			}
		} else {
			switch update.Message.Text {
			case "Курсы валют":
				mainCurrency := []string{"USD", "EUR"}
				data := exhanges{}

				res, err := http.Get(currencyURL)
				if err != nil {
					log.Println(err)
				}

				defer res.Body.Close()

				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Println(err)
				}

				err = json.Unmarshal(body, &data)
				if err != nil {
					log.Println(err)
				}

				msg.Text = getCurrencyValue(data, mainCurrency)

			case "Точное время (московское)":
				data := exactTime{}

				res, err := http.Get(timeURL)
				if err != nil {
					log.Println(err)
				}

				defer res.Body.Close()

				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Println(err)
				}

				err = json.Unmarshal(body, &data)
				if err != nil {
					log.Println(err)
				}

				msg.Text = getTimeAsString(data.Time)

			default:
				msg.Text = "I don't know what is this..."
			}
		}

		msg.ReplyMarkup = buttonsMarkup
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}
