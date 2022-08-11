package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// for override confirm
// "existed_filename.txt": "new content"
var Queued_Overrides = make(map[string]string)

func handleUpdateMessage(bot *tgbotapi.BotAPI, Message *tgbotapi.Message) {
	if Message.IsCommand() {
		// handle commands
		switch Message.Command() {
		case "echo":
			replyMsg := tgbotapi.NewMessage(Message.Chat.ID, Message.Text)
			replyMsg.ReplyToMessageID = Message.MessageID
			if _, err := bot.Send(replyMsg); err != nil {
				log.Println(err)
			}
		case "new", "add": // new hok tse bun
			// find file name
			split_tmp := strings.Split(Message.Text, " ")
			if len(split_tmp) <= 2 {
				replyMsg := tgbotapi.NewMessage(Message.Chat.ID, fmt.Sprintf("錯誤：新增格式爲 “/%s {關鍵字} {內容}”", Message.Command()))
				replyMsg.ReplyToMessageID = Message.MessageID
				if _, err := bot.Send(replyMsg); err != nil {
					log.Println(err)
				}
				return
			}

			// check file existence
			var filename string = split_tmp[1] + ".txt"
			var content string = Message.Text[len(Message.Command())+len(filename)-1:]
			content = strings.TrimSpace(content)
			if v, is_exist := CACHE[delExtension(filename)]; is_exist {
				if utf8.RuneCountInString(v) >= 100 {
					r := []rune(v)[:100]
					v = string(r) + "……"
				}
				Queued_Overrides[filename] = content
				replyMsg := tgbotapi.NewMessage(Message.Chat.ID, fmt.Sprintf("「%s」複製文已存在：「%s」，確認是否覆蓋？", split_tmp[1], v))
				replyMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("是", filename),
						tgbotapi.NewInlineKeyboardButtonData("否", "NIL"),
					),
				)
				if _, err := bot.Send(replyMsg); err != nil {
					log.Println(err)
				}
				return
			}
			// write file
			file, err := os.Create(path.Join(FILE_LOCATION, filename))
			if err != nil {
				log.Println(err)
			}
			file.WriteString(content)
			file.Close()

			// update cache
			CACHE[delExtension(filename)] = content

			// send response to user
			replyMsg := tgbotapi.NewMessage(Message.Chat.ID, fmt.Sprintf("新增複製文「%s」成功", delExtension(filename)))
			replyMsg.ReplyToMessageID = Message.MessageID
			if _, err := bot.Send(replyMsg); err != nil {
				log.Println(err)
			}
		case "random":
			k := rand.Intn(len(CACHE))
			var context string
			for _, v := range CACHE {
				if k == 0 {
					context = v
				}
				k--
			}
			context = fmt.Sprintf("幫你從 %d 篇文章中精心選擇了：\n%s", len(CACHE), context)
			replyMsg := tgbotapi.NewMessage(Message.Chat.ID, context)
			if _, err := bot.Send(replyMsg); err != nil {
				log.Println(err)
			}
		}
	} else {
		// search hok tse bun
		for k := range CACHE {
			if strings.Contains(Message.Text, k) {
				// hit
				replyMsg := tgbotapi.NewMessage(Message.Chat.ID, CACHE[k])
				if _, err := bot.Send(replyMsg); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func handleUpdateCallbackQuery(bot *tgbotapi.BotAPI, CallbackQuery *tgbotapi.CallbackQuery) {
	if CallbackQuery.Data == "NIL" {
		replyMsg := tgbotapi.NewMessage(CallbackQuery.Message.Chat.ID, "其實不按否也沒差啦 哈哈")
		replyMsg.ReplyToMessageID = CallbackQuery.Message.MessageID
		if _, err := bot.Send(replyMsg); err != nil {
			log.Println(err)
		}
	} else {
		var filename string = CallbackQuery.Data
		var content string = Queued_Overrides[filename]
		// write file
		file, err := os.Create(path.Join(FILE_LOCATION, filename))
		if err != nil {
			panic(err)
		}
		file.WriteString(content)
		file.Close()

		// update cache
		CACHE[delExtension(filename)] = content

		// send response to user
		replyMsg := tgbotapi.NewMessage(CallbackQuery.Message.Chat.ID, fmt.Sprintf("更新複製文「%s」成功", delExtension(filename)))
		replyMsg.ReplyToMessageID = CallbackQuery.Message.MessageID
		if _, err := bot.Send(replyMsg); err != nil {
			log.Println(err)
		}
	}
	editedMsg := tgbotapi.NewEditMessageReplyMarkup(
		CallbackQuery.Message.Chat.ID,
		CallbackQuery.Message.MessageID,
		tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: make([][]tgbotapi.InlineKeyboardButton, 0),
		},
	)
	bot.Send(editedMsg)
}

func main() {
	// keep alive
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "Hello World!")
	})
	go http.ListenAndServe(":9000", nil)

	// initialize
	file, _ := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	log.SetOutput(file)

	if _, err := os.Stat(FILE_LOCATION); os.IsNotExist(err) {
		os.Mkdir(FILE_LOCATION, 0755)
	}
	build_cache()

	// start bot
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	fmt.Println("***", "Sucessful logged in as", bot.Self.UserName, "***")

	// update config
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// get messages
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		// ignore nil
		if update.Message != nil {
			handleUpdateMessage(bot, update.Message)
		} else if update.CallbackQuery != nil {
			handleUpdateCallbackQuery(bot, update.CallbackQuery)
		}

	}

}
