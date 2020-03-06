package cmd

import (
	"fmt"

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

	// Validate we can find all users in slack
	slackUserTriggerrer, err := c.GetSlackUser(c.Config.Triggerrer)
	if err != nil {
		return 1, err
	}

	log.Debugf("found triggerrer slack user ID for %s: %s", c.Config.Triggerrer, slackUserTriggerrer.ID)

	slackUserReviewers := map[string]*slack.User{}
	for _, u := range c.Config.Reviewers {
		slackUser, err := c.GetSlackUser(u)
		if err != nil {
			return 1, err
		}
		slackUserReviewers[slackUser.ID] = slackUser
		log.Debugf("found reviewer slack user ID for %s : %s", u, slackUser.ID)
	}

	attachment := slack.Attachment{
		Title:      fmt.Sprintf("%s ( triggerred by @%s)", c.Config.Slack.Message, slackUserTriggerrer.ID),
		CallbackID: connectionID,
		Color:      "#3AA3E3",
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:  "approve",
				Text:  "Approve",
				Type:  "button",
				Style: "primary",
			},
			slack.AttachmentAction{
				Name:  "deny",
				Text:  "Deny",
				Type:  "button",
				Style: "danger",
			},
		},
	}

	channelID, messageTimestamp, err := c.Slack.PostMessage(c.Config.Slack.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		log.Fatal(err)
	}

	log.Debug(channelID)
	log.Debug(messageTimestamp)

	ok, err := c.ListenForApprovals(slackUserReviewers)
	if err != nil {
		return 1, err
	}

	if !ok {
		return 1, nil
	}

	return 0, nil
}
