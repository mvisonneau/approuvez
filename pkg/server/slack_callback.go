package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	pb "github.com/mvisonneau/approuvez/pkg/protobuf"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

const (
	approveMessage = ":white_check_mark: approved!"
	denyMessage    = ":x: denied!"
)

// HandleSlackCallback ..
func (s *Server) HandleSlackCallback(w http.ResponseWriter, r *http.Request) {
	var p slack.InteractionCallback
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r.Body); err != nil {
		log.WithField("error", err).Errorf("unable to read the payload's body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	unescapedBody, err := url.QueryUnescape(buf.String()[8:])
	if err != nil {
		log.WithField("error", err).Errorf("unable to unescape the body payload")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal([]byte(unescapedBody), &p); err != nil {
		log.WithField("payload", unescapedBody).Debugf("slack callback payload")
		log.WithField("error", err).Errorf("unable to parse JSON from slack callback payload")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Recompose our message from the callback payload
	msg := SlackMessage{}
	msg.Recompose(p.Message)
	msg.ActionButtons = false

	if len(p.ActionCallback.BlockActions) == 0 {
		log.Errorf("could not determine the action from the block")
		http.Error(w, "", http.StatusBadRequest)

		msg.StatusMessage = ":warning: an error occurred when sending the response (unknown action)"
		msg.Color = messageColorOrange
		s.UpdateMessage(p.Channel.ID, p.Message.Timestamp, msg)
		return
	}

	var decision pb.SlackUserResponse_Decision
	msg.StatusMessage = denyMessage
	msg.Color = messageColorRed
	action := p.ActionCallback.BlockActions[0]
	switch action.ActionID {
	case "linkButton":
		// When URL buttons are clicked we still receive an interaction payload
		// https://api.slack.com/reference/block-kit/block-elements#button__fields
		// In such cases we simply need to return a 200 to acknowledge.
		log.Debug("URL button being clicked, simply acknowledging the interaction request..")
		return
	case "approve":
		decision = pb.SlackUserResponse_APPROVE
		msg.StatusMessage = approveMessage
		msg.Color = messageColorGreen
	}

	sessionID, err := uuid.Parse(msg.SessionID)
	if err != nil {
		log.WithField("error", err).WithField("session_id", p.BlockID).Errorf("invalid session id from slack callback payload")
		http.Error(w, err.Error(), http.StatusBadRequest)
		msg.StatusMessage = ":warning: an error occurred when sending the response (invalid session_id)"
		msg.Color = messageColorOrange
		s.UpdateMessage(p.Channel.ID, p.Message.Timestamp, msg)
		return
	}

	session, ok := s.Sessions[sessionID.String()]
	if !ok {
		log.Errorf("session with the client does not exist")
		http.Error(w, "", http.StatusBadRequest)
		msg.StatusMessage = ":warning: an error occurred when sending the response (session closed)"
		msg.Color = messageColorOrange
		s.UpdateMessage(p.Channel.ID, p.Message.Timestamp, msg)
		return
	}

	log.WithFields(log.Fields{
		"session_id": sessionID.String(),
		"user_id":    p.User.ID,
		"decision":   decision,
	}).Info("processing slack callback")

	s.UpdateMessage(p.Channel.ID, p.Message.Timestamp, msg)

	err = session.Stream.Send(&pb.SlackUserResponse{
		User: &pb.SlackUser{
			Id:   p.User.ID,
			Name: p.User.Name,
		},
		Decision: decision,
	})

	if err != nil {
		log.WithField("error", err).Errorf("unable to forward the response to the client")
		delete(s.Sessions, session.ID)

		msg.StatusMessage = ":warning: the session with the client was already closed, response ignored"
		msg.Color = messageColorOrange
		s.UpdateMessage(p.Channel.ID, p.Message.Timestamp, msg)

		return
	}
}
