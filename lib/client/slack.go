package client

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
)

// GetSlackUser returns a slack user based on its email, name or ID
func (c *Client) GetSlackUser(ref string) (*slack.User, error) {
	// Validate we have passed an email
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(ref) {
		return c.Slack.GetUserByEmail(ref)
	}

	return c.Slack.GetUserInfo(ref)
}

// GenerateMessageBlocks compute the message blocks
func (c *Client) GenerateMessageBlocks(triggerrer *slack.User, reviewers map[string]*slack.User, decisions map[string]bool) []slack.Block {
	headerText := slack.NewTextBlockObject("mrkdwn", c.Config.Slack.Message, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	triggerrerText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Triggered by:* <@%s>", triggerrer.ID), false, false)
	triggerrerSection := slack.NewSectionBlock(triggerrerText, nil, nil)

	reviewersText := slack.NewTextBlockObject("mrkdwn", c.computeReviewersText(reviewers, decisions), false, false)
	reviewersSection := slack.NewSectionBlock(reviewersText, nil, nil)

	return []slack.Block{
		headerSection,
		triggerrerSection,
		reviewersSection,
	}
}

func (c *Client) computeReviewersText(reviewers map[string]*slack.User, decisions map[string]bool) (msg string) {
	remainers, approvers, deniers := []string{}, []string{}, []string{}

	// Find all the reviewers which have not replied yet
	for userID := range reviewers {
		if _, ok := decisions[userID]; !ok {
			remainers = append(remainers, fmt.Sprintf("<@%s>", userID))
		}
	}

	// Find all the reviewers which have approved or denied
	if len(decisions) > 0 {
		for userID, d := range decisions {
			if d {
				approvers = append(approvers, fmt.Sprintf("<@%s>", userID))
			} else {
				deniers = append(deniers, fmt.Sprintf("<@%s>", userID))
			}
		}

		if len(approvers) > 0 {
			msg += fmt.Sprintf("✅ approved by %s\n", strings.Join(approvers, ", "))
		}

		if len(deniers) > 0 {
			msg += fmt.Sprintf("❌ denied by %s", strings.Join(deniers, ", "))
			return
		}
	}

	if len(remainers) > 0 {
		msg += fmt.Sprintf("waiting to be approved by %s", strings.Join(remainers, ", "))
	}

	return
}

// SubmitCancellationMessages ..
func (c *Client) SubmitCancellationMessages(messages *Messages) error {
	// TODO: Do not replace this message entirely
	headerText := slack.NewTextBlockObject("mrkdwn", c.Config.Slack.Message, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	cancelText := slack.NewTextBlockObject("mrkdwn", "🐛 oops, this one got cancelled! 🤷‍♂️", false, false)
	cancelSection := slack.NewSectionBlock(cancelText, nil, nil)

	_, _, _, err := c.Slack.UpdateMessage(messages.Channel.ChannelID, messages.Channel.MessageTimestamp, slack.MsgOptionBlocks(headerSection, cancelSection))
	if err != nil {
		return err
	}

	for _, m := range messages.Users {
		if _, ok := m["action"]; ok {
			_, _, _, err = c.Slack.UpdateMessage(m["action"].ChannelID, m["action"].MessageTimestamp, slack.MsgOptionText("🐛 oops, this one got cancelled! 🤷‍♂️", false), slack.MsgOptionAttachments(slack.Attachment{}))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
