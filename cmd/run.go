package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mvisonneau/approuvez/lib/client"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

// Run is the main entrypoint of the tool
func Run(ctx *cli.Context) (int, error) {
	c, err := configure(ctx)
	if err != nil {
		return 1, err
	}
	defer c.Websocket.Close()

	// Fetch websocket connection ID
	connectionID, err := c.GetConnectionID()
	if err != nil {
		return 1, err
	}

	log.Debugf("Fetching users on Slack")
	triggerrer, err := c.GetSlackUser(c.Config.Triggerrer)
	if err != nil {
		return 1, fmt.Errorf("slack error: %v", err)
	}

	log.Debugf("Found triggerrer slack user ID for %s: %s", c.Config.Triggerrer, triggerrer.ID)

	reviewers := map[string]*slack.User{}
	for _, u := range c.Config.Reviewers {
		slackUser, err := c.GetSlackUser(u)
		if err != nil {
			return 1, fmt.Errorf("slack error: %v", err)
		}
		reviewers[slackUser.ID] = slackUser
		log.Debugf("Found reviewer slack user ID for %s : %s", u, slackUser.ID)
	}

	// Initialise a messages variable to store the references to every message we send
	messages := client.Messages{}
	messages.Users = map[string]map[string]*client.MessageRef{}

	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interruptChan
		c.SubmitCancellationMessages(&messages)
		log.Fatal("received interrupt, exiting with error 1")
	}()

	// Send Message
	channelID, messageTimestamp, err := c.Slack.PostMessage(c.Config.Slack.Channel, slack.MsgOptionBlocks(c.GenerateMessageBlocks(triggerrer, reviewers, map[string]bool{})...))
	if err != nil {
		if err := c.SubmitCancellationMessages(&messages); err != nil {
			return 1, err
		}
		return 1, err
	}

	messages.Channel = &client.MessageRef{
		ChannelID:        channelID,
		MessageTimestamp: messageTimestamp,
	}

	permalink, err := c.Slack.GetPermalink(&slack.PermalinkParameters{Channel: channelID, Ts: messageTimestamp})
	if err != nil {
		if err := c.SubmitCancellationMessages(&messages); err != nil {
			return 1, err
		}
		return 1, err
	}

	for userID := range reviewers {
		messages.Users[userID] = map[string]*client.MessageRef{}

		channelID, messageTimestamp, err := c.Slack.PostMessage(userID, slack.MsgOptionText(permalink, false))
		if err != nil {
			if err := c.SubmitCancellationMessages(&messages); err != nil {
				return 1, err
			}
			return 1, err
		}

		messages.Users[userID]["link"] = &client.MessageRef{
			ChannelID:        channelID,
			MessageTimestamp: messageTimestamp,
		}

		attachment := slack.Attachment{
			CallbackID: connectionID,
			Color:      "#3AA3E3",
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

		channelID, messageTimestamp, err = c.Slack.PostMessage(userID, slack.MsgOptionAttachments(attachment))
		if err != nil {
			if err := c.SubmitCancellationMessages(&messages); err != nil {
				return 1, err
			}
			return 1, err
		}

		messages.Users[userID]["action"] = &client.MessageRef{
			ChannelID:        channelID,
			MessageTimestamp: messageTimestamp,
		}
	}

	ok, err := c.ListenForApprovals(&messages, triggerrer, reviewers)
	if err != nil {
		if err := c.SubmitCancellationMessages(&messages); err != nil {
			return 1, err
		}
		return 1, err
	}

	if !ok {
		return 1, nil
	}

	return 0, nil
}
