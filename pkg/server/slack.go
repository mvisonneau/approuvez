package server

import (
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

const (
	messageColorDefault = "#6d87d1"
	messageColorGreen   = "#36a64f"
	messageColorRed     = "#cc1212"
	messageColorOrange  = "#e8971e"
)

// SlackMessage ..
type SlackMessage struct {
	SessionID      string
	Message        string
	LinkButtonName string
	LinkButtonURL  string
	ActionButtons  bool
	StatusMessage  string
	Color          string
}

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

// Render compute the message blocks
func (msg *SlackMessage) Render() slack.MsgOption {
	section := slack.NewSectionBlock(
		&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: strings.Replace(msg.Message, `\n`, "\n", -1),
		},
		nil,
		nil,
	)

	if len(msg.LinkButtonName) > 0 {
		linkButton := slack.NewButtonBlockElement(
			"linkButton",
			"",
			slack.NewTextBlockObject(
				slack.PlainTextType,
				msg.LinkButtonName,
				true,
				false,
			),
		)
		linkButton.URL = msg.LinkButtonURL
		section.Accessory = slack.NewAccessory(linkButton)
	}

	footer := slack.NewContextBlock(
		"footer",
		slack.NewTextBlockObject(
			slack.PlainTextType,
			msg.StatusMessage,
			true,
			false,
		),
	)

	approveButton := slack.ButtonBlockElement{
		Type:     slack.METButton,
		ActionID: "approve",
		Style:    slack.StylePrimary,
		Text: slack.NewTextBlockObject(
			slack.PlainTextType,
			"Approve",
			true,
			false,
		),
	}

	denyButton := slack.ButtonBlockElement{
		Type:     slack.METButton,
		ActionID: "deny",
		Style:    slack.StyleDanger,
		Text: slack.NewTextBlockObject(
			slack.PlainTextType,
			"Deny",
			true,
			false,
		),
	}

	buttons := slack.NewActionBlock(
		msg.SessionID,
		approveButton,
		denyButton,
	)

	blocks := []slack.Block{
		section,
		slack.NewDividerBlock(),
		footer,
	}

	if msg.ActionButtons {
		blocks = append(blocks, buttons)
	}

	return slack.MsgOptionAttachments(slack.Attachment{
		Color:  msg.Color,
		Blocks: slack.Blocks{BlockSet: blocks},
	})
}

// Recompose a message from an interaction payload
func (msg *SlackMessage) Recompose(payload slack.Message) {
	if len(payload.Attachments) != 1 {
		log.Error("expected the payload to contain exactly 1 attachement")
		return
	}

	msg.Color = payload.Attachments[0].Color
	newMessageBlocks := []slack.Block{}
	for _, b := range payload.Attachments[0].Blocks.BlockSet {
		switch b.BlockType() {
		case slack.MBTSection:
			msg.Message = b.(*slack.SectionBlock).Text.Text

			if b.(*slack.SectionBlock).Accessory != nil &&
				b.(*slack.SectionBlock).Accessory.ButtonElement != nil &&
				b.(*slack.SectionBlock).Accessory.ButtonElement.Text != nil {
				msg.LinkButtonName = b.(*slack.SectionBlock).Accessory.ButtonElement.Text.Text
				msg.LinkButtonURL = b.(*slack.SectionBlock).Accessory.ButtonElement.URL
			}
		case slack.MBTContext:
			if b.(*slack.ContextBlock).BlockID == "footer" && len(b.(*slack.ContextBlock).ContextElements.Elements) == 1 {
				msg.StatusMessage = string(b.(*slack.ContextBlock).ContextElements.Elements[0].(*slack.TextBlockObject).Text)
			}
		case slack.MBTAction:
			msg.ActionButtons = true
			msg.SessionID = b.(*slack.ActionBlock).BlockID
		default:
			newMessageBlocks = append(newMessageBlocks, b)
		}
	}
}

// PromptSlackUser ..
func (s *Server) PromptSlackUser(msg SlackMessage, userID string) error {
	log.WithFields(log.Fields{
		"session_id": msg.SessionID,
		"user_id":    userID,
	}).Debug("prompting slack user")

	_, _, err := s.Slack.PostMessage(userID, msg.Render())
	return err
}

// UpdateMessage ..
func (s *Server) UpdateMessage(channelID, messageTimestamp string, msg SlackMessage) {
	log.WithFields(log.Fields{
		"session_id": msg.SessionID,
	}).Debug("updating slack message")

	_, _, _, err := s.Slack.UpdateMessage(channelID, messageTimestamp, msg.Render())
	if err != nil {
		log.WithField("error", err).Errorf("unable to update the slack message")
	}
}
