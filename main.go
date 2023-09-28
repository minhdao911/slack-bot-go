package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/minhdao911/reminder-slack-bot-go/controllers"
	"github.com/minhdao911/reminder-slack-bot-go/drivers"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	godotenv.Load()

	client, err := drivers.ConnectToSlackViaSocketmode();
	if err != nil {
		log.Fatal(err)
	}
	
	socketmodeHandler := socketmode.NewSocketmodeHandler(client)

	controllers.NewSlashCommandController(socketmodeHandler)
	controllers.NewReminderController(socketmodeHandler)

	socketmodeHandler.RunEventLoop()
}
