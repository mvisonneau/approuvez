package client

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// GetSlackUser returns a slack user based on its email, name or ID
func (c *Client) GetSlackUser(ref string) (*slack.User, error) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(ref) {
		log.Debugf("Looking up Slack user '%s' using email", ref)
		return c.Slack.GetUserByEmail(ref)
	}

	log.Debugf("Looking up Slack user '%s' using ID", ref)
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

// SubmitApprovalMessages ..
func (c *Client) SubmitApprovalMessages(messages *Messages, triggerrer *slack.User, reviewers map[string]*slack.User, decisions map[string]bool, userID string) error {
	// Update the status of the message on Slack
	_, _, _, err := c.Slack.UpdateMessage(messages.Channel.ChannelID, messages.Channel.MessageTimestamp, slack.MsgOptionBlocks(c.GenerateMessageBlocks(triggerrer, reviewers, decisions)...))
	if err != nil {
		return err
	}

	// Replace the buttons
	_, _, _, err = c.Slack.UpdateMessage(messages.Users[userID]["action"].ChannelID, messages.Users[userID]["action"].MessageTimestamp, slack.MsgOptionText("you approved ✅ !", false), slack.MsgOptionAttachments(slack.Attachment{}))
	if err != nil {
		return err
	}

	return nil
}

// SubmitDenialMessages ..
func (c *Client) SubmitDenialMessages(messages *Messages, triggerrer *slack.User, reviewers map[string]*slack.User, decisions map[string]bool, userID string) error {
	// Update the status of the message on Slack
	_, _, _, err := c.Slack.UpdateMessage(messages.Channel.ChannelID, messages.Channel.MessageTimestamp, slack.MsgOptionBlocks(c.GenerateMessageBlocks(triggerrer, reviewers, decisions)...))
	if err != nil {
		return err
	}

	// Remove buttons for current users
	_, _, _, err = c.Slack.UpdateMessage(messages.Users[userID]["action"].ChannelID, messages.Users[userID]["action"].MessageTimestamp, slack.MsgOptionText("you denied ❌ !", false), slack.MsgOptionAttachments(slack.Attachment{}))
	if err != nil {
		return err
	}

	// Remove buttons for other reviewers
	for u, m := range messages.Users {
		if u != userID {
			_, _, _, err = c.Slack.UpdateMessage(m["action"].ChannelID, m["action"].MessageTimestamp, slack.MsgOptionText(fmt.Sprintf("denied by <@%s> ❌ !", userID), false), slack.MsgOptionAttachments(slack.Attachment{}))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
