package client

import (
	"strings"

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

// NewClient instantiate a Client from a provider Config
func NewClient(input *NewClientInput) (c *Client, err error) {
	var ws *websocket.Conn
	ws, _, err = websocket.DefaultDialer.Dial(input.WebsocketEndpoint, nil)
	if err != nil {
		return
	}

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
func (c *Client) ListenForApprovals(slackUserReviewers map[string]*slack.User) (bool, error) {
	requiredApprovals := len(slackUserReviewers)
	if c.Config.RequiredApprovals > 0 {
		requiredApprovals = c.Config.RequiredApprovals
	}

	for requiredApprovals > 0 {
		if requiredApprovals == 1 {
			log.Info("1 more approval required, waiting for it..")
		} else {
			log.Infof("%d more approvals required, waiting for it..", requiredApprovals)
		}

		_, resp, err := c.Websocket.ReadMessage()
		if err != nil {
			return false, err
		}

		s := strings.Split(string(resp), "/")
		if len(s) != 2 {
			log.Warnf("unable to interprete response '%s' 🤷‍♂️", string(resp))
			continue
		}

		userID := s[0]
		decision := s[1]
		log.Debugf("response - user id: %s / decision %s", userID, decision)

		if _, ok := slackUserReviewers[userID]; !ok {
			log.Warnf("Received a response from User ID '%s' but this user is not part for the allowed reviewers, skipping..", userID)
			continue
		}

		switch decision {
		case "approve":
			log.Infof("approved by %s!", slackUserReviewers[userID].Name)
			requiredApprovals--
		case "deny":
			log.Infof("denied by %s! exiting", slackUserReviewers[userID].Name)
			return false, nil
		default:
			log.Infof("unable to interprete decision '%s' 🤷‍♂️", decision)
		}
	}

	return true, nil
}
