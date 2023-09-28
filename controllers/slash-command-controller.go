package controllers

import (
	"log"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type SlashCommandController struct {
	EventHandler *socketmode.SocketmodeHandler
}

func NewSlashCommandController(handler *socketmode.SocketmodeHandler) {
	c := SlashCommandController{
		EventHandler: handler,
	}

	c.EventHandler.HandleSlashCommand("/ping", c.pong)
}

func (c SlashCommandController) pong(evt *socketmode.Event, clt *socketmode.Client) {
	command, ok := evt.Data.(slack.SlashCommand)
	if !ok {
		log.Printf("ERROR converting event to Slash Command: %v", command)
	}else{
		log.Printf("Slack Command received: %v", command.Command)
	}

	clt.Ack(*evt.Request)

	msg := slack.MsgOptionText("pong!", false)
	_, _, err := clt.PostMessage(command.ChannelID, msg)
	if err != nil {
		log.Printf("ERROR while sending message: %v", err)
	}
}
