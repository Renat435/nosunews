package main

import (
	"context"
	"fmt"
	"getNewsByDate/internal/utils"
	"github.com/PuerkitoBio/goquery"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func GetTitles(pageNumber int, c chan map[string]string, wg *sync.WaitGroup) {

	wg.Add(1)
	url := "https://www.nosu.ru/category/news/#"
	formData := make(map[string]string)
	formData["paged"] = fmt.Sprint(pageNumber)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	page, _, err := utils.GetPage(context.Background(), http.MethodPost, url, nil, headers, formData, 0)
	if err != nil {
		log.Println(err)
	}
	page.Find(".content-block").Each(func(i int, s *goquery.Selection) {

		imgPath, _ := s.Find("img").Attr("src")
		href, _ := s.Find(".title > a").Attr("href")

		result := map[string]string{
			"date":  s.Find(".date").Text(),
			"title": s.Find("a").Text(),
			"text":  s.Find("p").Text(),
			"href":  href,
			"img":   imgPath,
		}

		c <- result
	})
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan map[string]string, 32)
	for i := 1; i < 60; i++ {
		idx := i
		go GetTitles(idx, ch, &wg)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	// инициализация бота
	var bot, err = tgBotApi.NewBotAPI("6225547673:AAExptFAggw4gO6fucVpExovUWMFmVmpjcE")

	// остановка бота если не получилось подключиться
	if err != nil {
		log.Panic(err)
	}

	// включение отладки
	bot.Debug = true

	// выводим в консоль сообщение о том что мы успешно подключились
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// обновление конфигурации
	updateConfig := tgBotApi.NewUpdate(0)
	updateConfig.Timeout = 60

	// получение новых обновлений с новыми настройками
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		// если сообщение равно nil, то проходим дальше
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/start" {

			msg := tgBotApi.NewMessage(update.Message.Chat.ID, "Для того чтобы найти новость введите дату в формате дд.мм.гггг")

			numericKeyboard := tgBotApi.NewReplyKeyboard(
				tgBotApi.NewKeyboardButtonRow(
					tgBotApi.NewKeyboardButton("Сегодняшние новости"),
				),
			)

			msg.ReplyMarkup = numericKeyboard

			_, err := bot.Send(msg)

			if err != nil {
				panic(err)
			}
			continue
		}

		t, DateErr := time.Parse("02.01.2006", update.Message.Text)

		msg := tgBotApi.NewMessage(update.Message.Chat.ID, "Вы ввели неправильную дату \nВведите дату в формате дд.мм.гггг")

		if update.Message.Text == "Сегодняшние новости" {
			DateErr = nil
			t = time.Now()
		}

		if DateErr == nil {

			newsQnt := 0

			for r := range ch {
				if r["date"] == t.Format("02.01.2006") {

					photo := tgBotApi.NewPhoto(update.Message.From.ID, tgBotApi.FileURL(r["img"]))
					photo.Caption = r["title"] + "\n \n" + r["text"]

					var test = tgBotApi.NewInlineKeyboardMarkup(
						tgBotApi.NewInlineKeyboardRow(
							tgBotApi.NewInlineKeyboardButtonURL("Подробнее", r["href"]),
						),
					)

					photo.ReplyMarkup = test

					if _, err = bot.Send(photo); err != nil {
						log.Fatalln(err)
					}

					newsQnt++

				}
			}

			newsQntText := "В эту дату не было опубликованно ни одной новости"

			if update.Message.Text == "Сегодняшние новости" {
				newsQntText = "Сегодня ещё не было опубликованно новостей)"
			}

			if newsQnt > 0 {
				newsQntText = "Количество новостей опубликованных в эту дату: " + strconv.Itoa(newsQnt)
			}

			msg = tgBotApi.NewMessage(update.Message.Chat.ID, newsQntText)
		}

		_, err := bot.Send(msg)

		if err != nil {
			panic(err)
		}
	}
}
