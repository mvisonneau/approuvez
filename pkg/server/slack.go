package server

import (
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// GetSlackUser returns a slack user based on its email, name or ID
func (s *Server) GetSlackUser(ref string) (*slack.User, error) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(ref) {
		log.WithFields(log.Fields{
			"user-string": ref,
			"used-method": "email",
		}).Debug("looking up slack user")
		return s.Slack.GetUserByEmail(ref)
	}

	log.WithFields(log.Fields{
		"user-string": ref,
		"used-method": "id",
	}).Debug("looking up slack user")
	return s.Slack.GetUserInfo(ref)
}

// GenerateSlackMessageBlocks compute the message blocks
func GenerateSlackMessageBlocks(message, _ string) []slack.Block {
	headerText := slack.NewTextBlockObject("mrkdwn", message, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// reviewersText := slack.NewTextBlockObject("mrkdwn", user, false, false)
	// reviewersSection := slack.NewSectionBlock(reviewersText, nil, nil)

	return []slack.Block{
		headerSection,
		// reviewersSection,
	}
}

// PromptSlackUser ..
func (s *Server) PromptSlackUser(sessionID, message, userID string) error {
	log.WithFields(log.Fields{
		"user_id": userID,
	}).Debug("prompting slack user")

	attachment := slack.Attachment{
		CallbackID: sessionID,
		Color:      "#3AA3E3",
		Text:       "👋 what do you reckon?",
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

	_, _, err := s.Slack.PostMessage(userID, slack.MsgOptionText(message, false), slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("posting notification message: %v", err)
	}

	// channelID, messageTimestamp, err := c.Slack.PostMessage(userID, slack.MsgOptionText(permalink, false))
	// 	if err != nil {
	// 		log.WithFields(
	// 			log.Fields{
	// 				"user-id": userID,
	// 				"error":   err.Error(),
	// 			},
	// 		).Error("posting message link to user")

	// 		if err := c.SubmitCancellationMessages(messages); err != nil {
	// 			return 1, err
	// 		}
	// 		return 1, err
	// 	}

	// 	messages.Users[userID]["link"] = client.MessageRef{
	// 		ChannelID:        channelID,
	// 		MessageTimestamp: messageTimestamp,
	// 	}

	// 	log.WithFields(
	// 		log.Fields{
	// 			"user-id": userID,
	// 		},
	// 	).Debug("posting actions message to user")
	// 	channelID, messageTimestamp, err = c.Slack.PostMessage(userID, slack.MsgOptionAttachments(attachment))
	// 	if err != nil {
	// 		log.WithFields(
	// 			log.Fields{
	// 				"user-id": userID,
	// 				"error":   err.Error(),
	// 			},
	// 		).Error("posting actions message to user")
	// 		if err := c.SubmitCancellationMessages(messages); err != nil {
	// 			return 1, err
	// 		}
	// 		return 1, err
	// 	}

	// 	messages.Users[userID]["action"] = client.MessageRef{
	// 		ChannelID:        channelID,
	// 		MessageTimestamp: messageTimestamp,
	// 	}
	// }
	return nil
}
