package controllers

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type ReminderController struct {
	EventHandler *socketmode.SocketmodeHandler
}

func NewReminderController(handler *socketmode.SocketmodeHandler) ReminderController {
	c := ReminderController{
		EventHandler: handler,
	}

	c.EventHandler.HandleEvents(slackevents.AppMention, c.setupReminder)

	return c
}

func (c *ReminderController) setupReminder(evt *socketmode.Event, clt *socketmode.Client) {
	evt_api, _ := evt.Data.(slackevents.EventsAPIEvent)
	evt_mentioned, ok := evt_api.InnerEvent.Data.(*slackevents.AppMentionEvent)

	if !ok {
		log.Printf("ERROR converting event to slackevents.AppMentionEvent: %v", evt_api)
	}else{
		log.Printf("Event received: %v", evt_mentioned)
	}

	clt.Ack(*evt.Request)

	separator := "-"
	reminder_text := evt_mentioned.Text
	msg := slack.MsgOptionText("Wrong syntax, please try again", false)
	if !strings.Contains(reminder_text, separator) {
		clt.PostEphemeral(
			evt_mentioned.Channel,
			evt_mentioned.User,
			msg,
		)
		return
	}

	texts := strings.Split(reminder_text, separator)

	// loc, _ := time.LoadLocation("Local")
	// layout := "Jan 2, 2006 at 3:04pm"
	reminder_time := strings.TrimSpace(texts[1])
	// unix_time , err := time.ParseInLocation(layout, reminder_time, loc)
	// if err != nil {
	// 	log.Printf("ERROR parsing time: %v\n", err)
	// }
	// reminder_time = strconv.FormatInt(unix_time.Unix(), 10)
	reminder_time = getUnixTimeString(reminder_time)
	log.Printf("Unix time: %v\n", reminder_time)

	contents := strings.Split(texts[0], " ")
	reminder_content := strings.Join(contents[1:], " ")

	// msg = slack.MsgOptionText(
	// 	fmt.Sprintf("Reminder added for %v - %v", reminder_content, reminder_time), 
	// 	false,
	// )
	// clt.PostEphemeral(
	// 	evt_mentioned.Channel,
	// 	evt_mentioned.User,
	// 	msg,
	// )

	_, _ , err := clt.ScheduleMessage(evt_mentioned.Channel, reminder_time, slack.MsgOptionText(
		fmt.Sprintf("Reminder: %v", reminder_content),
		false,
	))
	if err != nil {
		log.Printf("ERROR while scheduling message: %v", err)
	}
}

func getUnixTimeString (time_string string) string {
	var unix int64
	now := time.Now()
	if strings.Contains(time_string, "in") {
		re := regexp.MustCompile(`\d+`)
		d := re.FindAllString(time_string, -1)[0]
		dInt, _ := strconv.Atoi(d)
		if strings.Contains(time_string, "min") || strings.Contains(time_string, "mins"){
			unix = now.Add(time.Duration(dInt) * time.Minute).Unix()

		}else if strings.Contains(time_string, "hr") || strings.Contains(time_string, "hrs"){
			unix = now.Add(time.Duration(dInt) * time.Hour).Unix()
		}
	} else if strings.Contains(time_string, "tomorrow") {
		hr_min := strings.TrimSpace(strings.Split(time_string, "at")[1])
		hr := strings.Split(hr_min, ":")[0]
		hrInt, _ := strconv.Atoi(hr)
		min := strings.Split(hr_min, ":")[1]
		minInt, _ := strconv.Atoi(min)

		unix = time.Date(now.Year(), now.Month(), now.Day()+1, hrInt, minInt,  0, 0, now.Location()).Unix()
	} else {
		hr_min := strings.Split(time_string, ":")
		hrInt, _ := strconv.Atoi(hr_min[0])
		minInt, _ := strconv.Atoi(hr_min[1])

		unix = time.Date(now.Year(), now.Month(), now.Day(), hrInt, minInt, 0, 0, now.Location()).Unix()
	}

	return strconv.FormatInt(unix, 10)
}
