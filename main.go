package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv" // don't forget to install modules
	"githhub.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func main(){

	godotenv.Load(".env")

	toke := os.Getenv("SLACK_AUTH_TOKEN")
	chanId := os.Getenv("SLACK_CHANNEL_ID")

	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken)) // added op level app token

	/* 
	adds websocket to slack,
	adds debbuging,
	adds stdout log reporting,
	*/
	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	ctx, cancel := context.WithCancel(context.Background()) // added ctx var

	defer cancel()

	/*
	added for loop to check for new events,
	added Client.Run() call,
	created events chan,
	*/
	go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
		for {
			select {
			case <-ctx.Done():
				log.Println("Shutting down socket mode listener")
				return
			case event := <-socketClient.Events:
				//socketClient event listener
				switch event.Type {

				case socketmode.EventTypeEventsAPI:

					eventsAPI, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("Could not type cast the event to the Events API: %v\n", event)
						continue
					}
					socketClient.Ack(*event.Request)

					err := HandleEventMessage(eventsAPI, client) // updated error handling with event handler
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}

	}

	/*
	function to handle event message,
	func continues type switching.
	*/
	func HandleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client) error {
		switch event.Type {

		case slackevents.CallbackEvent:

			innerEvent := event.innerEvent
			
			switch evnt := innerEvent.Data.(type) {
				err := HandleAppMentionEventToBot(evnt, client)
				if err != nil {
					return err
				}
			}
		default:
			return errors.New("unsupoorted event type")
		}
		return nil
		
	func HandleAppMentionEventToBot(event *slackevents.AppMentionEvent, client *slack.Client) error {
		user, err := client.GetUserInfo(event.User)
			if err != nil {
				return err
			}
		 
		text := strings.ToLower(event.Text)
		 
		attachment := slack.Attachment{}
		 
			if strings.Contains(text, "hello") || strings.Contains(text, "hi") {
				attachment.Text = fmt.Sprintf("Greetings! %s", user.Name)
				attachment.Color = "#E01E5A"
			} else if strings.Contains(text, "weather") {
				attachment.Text = fmt.Sprintf("Weather is sunny today. %s", user.Name)
				attachment.Color = "#ECB22E"
			} else if strings.Contains(text, "schedule") {
				attachment.Text = fmt.Sprintf("Here's the schedule!(there would be a schedule here if I had one). %s", user.Name)
				attachment.Color = "#ECB22E"
			} else if strings.Contains(text, "salaries") {
				attachment.Text = fmt.Sprintf("Here's the SALARIES???!?!?!!(running out of ideas). %s", user.Name)
				attachment.Color = "#ECB22E"
			} else {
				attachment.Text = fmt.Sprintf("I am good. How are you %s?", user.Name)
				attachment.Color = "#2EB67D"
			}
			_, _, err = client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
			if err != nil {
				return fmt.Errorf("failed to post message: %w", err)
			}
			return nil
		}


	}