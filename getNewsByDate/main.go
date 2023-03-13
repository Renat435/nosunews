package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	url := "https://www.nosu.ru/category/news/"
	//request method
	method := "POST"

	var jsonData = []byte(`{
		"paged": "2",
	}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	// set HTTP request header Content-Type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

	// инициализация бота
	//bot, err := tgBotApi.NewBotAPI("6225547673:AAExptFAggw4gO6fucVpExovUWMFmVmpjcE")
	//
	//// остановка бота если не получилось подключиться
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//// включение отладки
	//bot.Debug = true
	//
	//// выводим в консоль сообщение о том что мы успешно подключились
	//log.Printf("Authorized on account %s", bot.Self.UserName)
	//
	//// обновление конфигурации
	//updateConfig := tgBotApi.NewUpdate(0)
	//updateConfig.Timeout = 60
	//
	//// получение новых обновлений с новыми настройками
	//updates := bot.GetUpdatesChan(updateConfig)
	//
	//for update := range updates {
	//
	//	// если сообщение равно nil, то проходим дальше
	//	if update.Message == nil {
	//		continue
	//	}
	//
	//	if update.Message.Text == "/start" {
	//		msg := tgBotApi.NewMessage(update.Message.Chat.ID, "Для того чтобы найти новость введите дату в формате дд.мм.гггг")
	//		_, err := bot.Send(msg)
	//
	//		if err != nil {
	//			panic(err)
	//		}
	//		continue
	//	}
	//
	//	t, DateErr := time.Parse("02.01.2006", update.Message.Text)
	//
	//	msg := tgBotApi.NewMessage(update.Message.Chat.ID, "вы ввели неправильную дату")
	//
	//	if DateErr == nil {
	//		//t.Format("02.01.2006")
	//		msg = tgBotApi.NewMessage(update.Message.Chat.ID, t.Format("02.01.2006"))
	//	}
	//
	//	//switch update.Message.Text {
	//	//case "/start":
	//	//	msg = tgBotApi.NewMessage(update.Message.Chat.ID, "Для того чтобы найти новость введите дату в формате дд.мм.гггг")
	//	//default:
	//	//	msg = tgBotApi.NewMessage(update.Message.Chat.ID, "Для того чтобы найти новость введите дату в формате дд.мм.гггг")
	//	//	//msg := tgBotApi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	//	//}
	//
	//	_, err := bot.Send(msg)
	//
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	//// создаём новое сообщение имея предыдущее ( из предыдущего самое главное достать id чата, второй параметр отвечает за текст )
	//	//
	//	//msg := tgBotApi.NewMessage(update.Message.Chat.ID, update.Message.Text+"sosi "+strconv.Itoa(update.Message.From.ID))
	//	//
	//	///*
	//	//	если мы хотим ответить на сообщение,
	//	//	то в ReplyToMessageID нужно записать id
	//	//	сообщения на который мы хотим ответить если
	//	//	ничего не писать, то будет бот отправит
	//	//	сообщение не отвечая на предыдущее
	//	//*/
	//	//msg.ReplyToMessageID = update.Message.MessageID
	//	//
	//	//// отправляем созданное нами сообщение и отлавливаем ошибку если есть ошибка выводим её
	//	//_, err := bot.Send(msg)
	//	//
	//	//if err != nil {
	//	//	panic(err)
	//	//}
	//}
}
