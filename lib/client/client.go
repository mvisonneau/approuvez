package client

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

// NewClientInput is a helper for the NewClient() function
type NewClientInput struct {
	SlackToken        string
	SlackChannel      string
	SlackMessage      string
	WebsocketEndpoint string
	Triggerrer        string
	Reviewers         []string
	RequiredApprovals int
}

// Client handles necessary components to run the app
type Client struct {
	Slack     *slack.Client
	Websocket *websocket.Conn

	Config struct {
		Slack struct {
			Channel string
			Message string
		}

		Triggerrer        string
		Reviewers         []string
		RequiredApprovals int
	}
}

// Messages ..
type Messages struct {
	Channel *MessageRef
	Users   map[string]map[string]*MessageRef
}

// MessageRef can be used to store references to sent messages
type MessageRef struct {
	ChannelID        string
	MessageTimestamp string
}

// NewClient instantiate a Client from a provider Config
func NewClient(input *NewClientInput) (c *Client, err error) {
	var ws *websocket.Conn
	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}

	log.Debugf("Connecting to websocket %s", input.WebsocketEndpoint)
	ws, _, err = dialer.Dial(input.WebsocketEndpoint, nil)
	if err != nil {
		return
	}
	log.Infof("Successfully connected to websocket endpoint")

	c = &Client{
		Slack:     slack.New(input.SlackToken),
		Websocket: ws,
	}

	c.Config.Slack.Channel = input.SlackChannel
	c.Config.Slack.Message = input.SlackMessage
	c.Config.Triggerrer = input.Triggerrer
	c.Config.Reviewers = input.Reviewers
	c.Config.RequiredApprovals = input.RequiredApprovals

	return
}

// ListenForApprovals will loop until the amount of required approvals is reached
// or if a deny event is received
func (c *Client) ListenForApprovals(messages *Messages, triggerrer *slack.User, reviewers map[string]*slack.User) (bool, error) {
	decisions := map[string]bool{}
	requiredApprovals := c.getRequiredApprovals(reviewers)

	for requiredApprovals > 0 {
		if requiredApprovals == 1 {
			log.Info("1 more approval required, waiting for it..")
		} else {
			log.Infof("%d more approvals required, waiting for it..", requiredApprovals)
		}

		userID, decision, err := c.readResponse()
		if err != nil {
			return false, err
		}

		if _, ok := reviewers[userID]; !ok {
			log.Warnf("Received a response from User ID '%s' but this user is not part for the allowed reviewers, skipping..", userID)
			continue
		}

		switch decision {
		case "approve":
			log.Infof("approved by %s!", reviewers[userID].Name)
			decisions[userID] = true
			requiredApprovals--
			if err := c.SubmitApprovalMessages(messages, triggerrer, reviewers, decisions, userID); err != nil {
				return false, err
			}

		case "deny":
			log.Infof("denied by %s! exiting", reviewers[userID].Name)
			decisions[userID] = false

			if err := c.SubmitDenialMessages(messages, triggerrer, reviewers, decisions, userID); err != nil {
				return false, err
			}

			return false, nil
		default:
			log.Infof("unable to interprete decision '%s' 🤷‍♂️", decision)
		}
	}

	// Remove buttons for people who did not reply
	for u, m := range messages.Users {
		if _, ok := decisions[u]; !ok {
			_, _, _, err := c.Slack.UpdateMessage(m["action"].ChannelID, m["action"].MessageTimestamp, slack.MsgOptionText("already approved ✅, sorry for the noise 🙇‍♂️!", false), slack.MsgOptionAttachments(slack.Attachment{}))
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

func (c *Client) getRequiredApprovals(reviewers map[string]*slack.User) (requiredApprovals int) {
	requiredApprovals = len(reviewers)
	if c.Config.RequiredApprovals > 0 {
		requiredApprovals = c.Config.RequiredApprovals
	}
	return
}

func (c *Client) readResponse() (userID string, decision string, err error) {
	_, resp, err := c.Websocket.ReadMessage()
	if err != nil {
		return "", "", err
	}

	s := strings.Split(string(resp), "/")
	if len(s) != 2 {
		return "", "", fmt.Errorf("unable to interprete response '%s' 🤷‍♂️", string(resp))
	}

	userID = s[0]
	decision = s[1]
	log.Debugf("response - user id: %s / decision %s", userID, decision)
	return
}
