package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func BotToken() string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Errorf("erorr with get token. Error: %s", err)
	}
	return os.Getenv("BOT_TOKEN")
}
