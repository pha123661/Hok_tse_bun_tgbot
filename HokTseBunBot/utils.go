package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	toml "github.com/pelletier/go-toml/v2"
)

var CONFIG cfg

type Dict map[string]interface{}
type Empty struct{}

type cfg struct {
	SETTING struct {
		TYPE struct {
			TXT, IMG, ANI, VID int
		}
		NAME struct {
			TXT, IMG, ANI, VID string
		}
		CONCURRENT struct {
			SUM, CAP struct {
				COOLDOWN int // ms
				LIMIT    int
			}
		}
		LOG_FILE        string
		EXAMPLE_PIC_DIR string
	}

	API struct {
		TG struct {
			TOKEN string
		}
		HF struct {
			TOKENs              []string
			CURRENT_TOKEN       string
			SUM_MODEL, MT_MODEL string
		}
	}

	DB struct {
		DIR, EXPORT_DIR, CFormat string
	}
}

// Gets collection name of given ChatID
func (Config cfg) GetColbyChatID(ChatID int64) string {
	return fmt.Sprintf(CONFIG.DB.CFormat, ChatID)
}

// Gets Chinese name of given Type
func (Config cfg) GetNameByType(Type int) string {
	switch Type {
	case Config.SETTING.TYPE.TXT:
		return Config.SETTING.NAME.TXT
	case Config.SETTING.TYPE.IMG:
		return Config.SETTING.NAME.IMG
	case Config.SETTING.TYPE.ANI:
		return Config.SETTING.NAME.ANI
	case Config.SETTING.TYPE.VID:
		return Config.SETTING.NAME.VID
	default:
		return "大便"
	}
}

func InitConfig(CONFIG_PATH string) {
	// parse toml file
	tomldata, err := os.ReadFile(CONFIG_PATH)
	if err != nil {
		log.Panicln("[InitConfig]", err)
	}

	if err := toml.Unmarshal(tomldata, &CONFIG); err != nil {
		log.Panicln("[InitConfig]", err)
	}

	// get secret configs from environment variables
	CONFIG.API.HF.TOKENs = strings.Fields(os.Getenv("API.HF.TOKENs"))
	CONFIG.API.TG.TOKEN = os.Getenv("API.TG.TOKEN")

	SetHFAPI()

	fmt.Println("********************\nConfig Loaded:")
	PrintStructAsTOML(CONFIG)
	fmt.Println("********************")

	var CreateDirIfNotExist = func(path string) {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			errr := os.Mkdir(path, 0755)
			if errr != nil {
				log.Panicln("[InitConfig]", errr)
			}
		}
	}

	CreateDirIfNotExist(CONFIG.DB.EXPORT_DIR)
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

func PrintStructAsTOML(v interface{}) error {
	buf := bytes.Buffer{}
	enc := toml.NewEncoder(&buf)
	enc.SetIndentTables(true)
	if err := enc.Encode(v); err != nil {
		return err
	}
	fmt.Println(buf.String())
	return nil
}

// Sends text message, set ReplyMsgID=0 to disable reply
func SendText(ChatID int64, Content string, ReplyMsgID int) tgbotapi.Message {
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

// Sends media message
func SendMultiMedia(ChatID int64, Caption string, FileID_Str string, Type int) *tgbotapi.APIResponse {
	var Msg *tgbotapi.APIResponse
	var err error
	FileID := tgbotapi.FileID(FileID_Str)
	switch Type {
	case 1:
		log.Println("Sending text by SendMultiMedia")
		return nil

	case 2:
		Config := tgbotapi.NewPhoto(ChatID, FileID)
		if Caption != "" {
			Config.Caption = Caption
		}
		Msg, err = bot.Request(Config)
		if err != nil {
			log.Println("[SendIR]", err)
		}

	case 3:
		Config := tgbotapi.NewAnimation(ChatID, FileID)
		if Caption != "" {
			Config.Caption = Caption
		}
		Msg, err = bot.Request(Config)
		if err != nil {
			log.Println("[SendIR]", err)
		}

	case 4:
		Config := tgbotapi.NewVideo(ChatID, FileID)
		if Caption != "" {
			Config.Caption = Caption
		}
		Msg, err = bot.Request(Config)
		if err != nil {
			log.Println("[SendIR]", err)
		}
	}
	return Msg
}
