package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mvisonneau/approuvez/pkg/client"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
)

// Run is the main entrypoint of the tool
func Run(ctx *cli.Context) (int, error) {
	c, err := configure(ctx)
	if err != nil {
		return 1, err
	}
	defer c.Websocket.Close()

	// Fetch websocket connection ID
	connectionID, err := client.GetConnectionID(c.Websocket)
	if err != nil {
		return 1, err
	}

	// Keep the websocket connection alive
	client.KeepAlive(c.Websocket, 10*time.Minute)

	log.Debug("fetching Slack users")
	triggerrer, err := c.GetSlackUser(c.Config.Triggerrer)
	if err != nil {
		return 1, fmt.Errorf("slack error: %v", err)
	}

	log.WithFields(
		log.Fields{
			"user-string": c.Config.Triggerrer,
			"user-id":     triggerrer.ID,
			"user-name":   triggerrer.Name,
			"user-type":   "triggerrer",
		},
	).Debug("found slack user ID")

	reviewers := map[string]*slack.User{}
	for _, u := range c.Config.Reviewers {
		slackUser, err := c.GetSlackUser(u)
		if err != nil {
			return 1, fmt.Errorf("slack error: %v", err)
		}
		reviewers[slackUser.ID] = slackUser
		log.WithFields(
			log.Fields{
				"user-string": u,
				"user-id":     slackUser.ID,
				"user-name":   slackUser.Name,
				"user-type":   "reviewer",
			},
		).Debug("found slack user ID")
	}

	// Initialise a messages variable to store the references to every message we send
	messages := &client.Messages{}
	messages.Users = map[string]map[string]client.MessageRef{}

	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interruptChan
		if err = c.SubmitCancellationMessages(messages); err != nil {
			log.Error(err.Error())
		}
		log.Fatal("received interrupt, exiting with error 1")
	}()

	// Send Message
	log.WithFields(
		log.Fields{
			"channel": c.Config.Slack.Channel,
		},
	).Debug("posting notification message")
	channelID, messageTimestamp, err := c.Slack.PostMessage(c.Config.Slack.Channel, slack.MsgOptionBlocks(c.GenerateMessageBlocks(triggerrer, reviewers, map[string]bool{})...))
	if err != nil {
		log.WithFields(
			log.Fields{
				"channel": c.Config.Slack.Channel,
				"error":   err.Error(),
			},
		).Error("posting notification message")

		if err := c.SubmitCancellationMessages(messages); err != nil {
			return 1, err
		}
		return 1, err
	}

	messages.Channel = client.MessageRef{
		ChannelID:        channelID,
		MessageTimestamp: messageTimestamp,
	}

	log.WithFields(
		log.Fields{
			"channel":           channelID,
			"message-timestamp": messageTimestamp,
		},
	).Debug("getting permalink")
	permalink, err := c.Slack.GetPermalink(&slack.PermalinkParameters{Channel: channelID, Ts: messageTimestamp})
	if err != nil {
		log.WithFields(
			log.Fields{
				"channel":           channelID,
				"message-timestamp": messageTimestamp,
				"error":             err.Error(),
			},
		).Error("getting permalink")

		if err := c.SubmitCancellationMessages(messages); err != nil {
			return 1, err
		}
		return 1, err
	}

	log.WithFields(
		log.Fields{
			"permalink": permalink,
		},
	).Debug("found permalink")

	for userID := range reviewers {
		log.WithFields(
			log.Fields{
				"user-id": userID,
			},
		).Debug("posting message link to user")

		messages.Users[userID] = map[string]client.MessageRef{}
		channelID, messageTimestamp, err := c.Slack.PostMessage(userID, slack.MsgOptionText(permalink, false))
		if err != nil {
			log.WithFields(
				log.Fields{
					"user-id": userID,
					"error":   err.Error(),
				},
			).Error("posting message link to user")

			if err := c.SubmitCancellationMessages(messages); err != nil {
				return 1, err
			}
			return 1, err
		}

		messages.Users[userID]["link"] = client.MessageRef{
			ChannelID:        channelID,
			MessageTimestamp: messageTimestamp,
		}

		attachment := slack.Attachment{
			CallbackID: connectionID,
			Color:      "#3AA3E3",
			Text:       "what do you reckon?",
			Actions: []slack.AttachmentAction{
				{
					Name:  "approve",
					Text:  "Approve",
					Type:  "button",
					Style: "primary",
				},
				{
					Name:  "deny",
					Text:  "Deny",
					Type:  "button",
					Style: "danger",
				},
			},
		}

		log.WithFields(
			log.Fields{
				"user-id": userID,
			},
		).Debug("posting actions message to user")
		channelID, messageTimestamp, err = c.Slack.PostMessage(userID, slack.MsgOptionAttachments(attachment))
		if err != nil {
			log.WithFields(
				log.Fields{
					"user-id": userID,
					"error":   err.Error(),
				},
			).Error("posting actions message to user")
			if err := c.SubmitCancellationMessages(messages); err != nil {
				return 1, err
			}
			return 1, err
		}

		messages.Users[userID]["action"] = client.MessageRef{
			ChannelID:        channelID,
			MessageTimestamp: messageTimestamp,
		}
	}

	log.Debug("listening for Slack responses")
	ok, err := c.ListenForApprovals(messages, triggerrer, reviewers)
	if err != nil {
		log.WithFields(
			log.Fields{
				"error": err.Error(),
			},
		).Error("listening for Slack responses")
		if err := c.SubmitCancellationMessages(messages); err != nil {
			return 1, err
		}
		return 1, err
	}

	if !ok {
		return 1, nil
	}

	return 0, nil
}
