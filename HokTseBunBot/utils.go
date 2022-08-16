package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	toml "github.com/BurntSushi/toml"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var CONFIG Config_Type

type Config_Type struct {
	SETTING struct {
		LOG_FILE string
	}

	API struct {
		TG struct {
			TOKEN string
		}
		HF struct {
			TOKENs        []string
			CURRENT_TOKEN string
			MODEL         string
		}
	}

	DB struct {
		DIR        string
		COLLECTION string
	}
}

func InitConfig(CONFIG_PATH string) {
	// parse toml file
	tomldata, err := os.ReadFile(CONFIG_PATH)
	if err != nil {
		log.Panicln("[InitConfig]", err)
	}
	if _, err := toml.Decode(string(tomldata), &CONFIG); err != nil {
		log.Panicln("[InitConfig]", err)
	}

	buf := new(bytes.Buffer)
	toml.NewEncoder(buf).Encode(CONFIG)
	fmt.Printf("********************\nConfig Loaded:\n%s\n********************\n", buf.String())

	// var CreateDirIfNotExist = func(path string) {
	// 	if _, err := os.Stat(path); os.IsNotExist(err) {
	// 		errr := os.Mkdir(path, 0755)
	// 		if errr != nil {
	// 			log.Panicln("[InitConfig]", errr)
	// 		}
	// 	}
	// }

	// CreateDirIfNotExist(CONFIG.DB.DIR)
}

func TruncateString(text string, width int) string {
	text = strings.TrimSpace(text)
	width = width - utf8.RuneCountInString("……")
	if utf8.RuneCountInString(text) > width {
		r := []rune(text)[:width]
		text = string(r) + "……"
	}
	return text
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// helper functions
func SendTextResult(ChatID int64, Content string, ReplyMsgID int) tgbotapi.Message {
	replyMsg := tgbotapi.NewMessage(ChatID, Content)
	if ReplyMsgID != 0 {
		replyMsg.ReplyToMessageID = ReplyMsgID
	}
	Msg, err := bot.Send(replyMsg)
	if err != nil {
		log.Println("[SendTR]", err)
	}
	return Msg
}

func SendImageResult(ChatID int64, Caption string, Content string) *tgbotapi.APIResponse {
	FileID := tgbotapi.FileID(Content)
	PhotoConfig := tgbotapi.NewPhoto(ChatID, FileID)
	PhotoConfig.Caption = Caption
	Msg, err := bot.Request(PhotoConfig)
	if err != nil {
		log.Println("[SendIR]", err)
	}
	return Msg
}