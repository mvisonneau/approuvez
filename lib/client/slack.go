package client

import (
	"regexp"

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

func generateMessageBlocks() []slack.Block {
	// Header Section
	headerText := slack.NewTextBlockObject("mrkdwn", "Releasing foo app into dev", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Fields
	typeField := slack.NewTextBlockObject("mrkdwn", "*Triggered by:* @max", false, false)
	fieldSlice := make([]*slack.TextBlockObject, 0)
	fieldSlice = append(fieldSlice, typeField)

	fieldsSection := slack.NewSectionBlock(nil, fieldSlice, nil)

	approveText := slack.NewTextBlockObject("mrkdwn", "*Waiting to be approved by @max", false, false)
	approveSection := slack.NewSectionBlock(approveText, nil, nil)

	return []slack.Block{
		headerSection,
		fieldsSection,
		approveSection,
	}
}

func getEphemeralMessageActionBlock(permalink string) *slack.ActionBlock {
	approveOrDenyBtn := slack.NewButtonBlockElement("", "approve", slack.NewTextBlockObject("plain_text", "Approve or Deny", false, false))
	approveOrDenyBtn.URL = permalink
	discardBtn := slack.NewButtonBlockElement("", "discard", slack.NewTextBlockObject("plain_text", "Discard", false, false))
	return slack.NewActionBlock("", approveOrDenyBtn, discardBtn)
}
