package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"reminder/internal/bot"
	"reminder/internal/database/mysql"
)



func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	token := os.Getenv("BOT_TOKEN")

	mysql.InitDB()
	bot.StartBot(token)
}
