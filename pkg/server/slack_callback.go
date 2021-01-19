package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	pb "github.com/mvisonneau/approuvez/pkg/protobuf/approuvez"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// HandleSlackCallback ..
func (s *Server) HandleSlackCallback(w http.ResponseWriter, r *http.Request) {
	var p slack.InteractionCallback
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

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

	sessionID, err := uuid.Parse(p.CallbackID)
	if err != nil {
		log.WithField("error", err).WithField("callback_id", p.CallbackID).Errorf("invalid callback id from slack callback payload")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, ok := s.Sessions[sessionID.String()]
	if !ok {
		log.Errorf("session with the client does not exist")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if len(p.ActionCallback.AttachmentActions) < 1 {
		log.Errorf("unable to read the decision from the payload")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var decision pb.SlackUserResponse_Decision
	switch p.ActionCallback.AttachmentActions[0].Name {
	case "approve":
		decision = pb.SlackUserResponse_APPROVE
	}

	log.WithFields(log.Fields{
		"session_id": sessionID.String(),
		"user_id":    p.User.ID,
		"decision":   decision,
	}).Info("processing slack callback")

	_, _, _, err = s.Slack.UpdateMessage(p.Channel.ID, p.MessageTs, slack.MsgOptionAttachments(slack.Attachment{}), slack.MsgOptionText(fmt.Sprintf("%s\n%s", p.OriginalMessage.Text, "approved! ✅"), false))
	if err != nil {
		log.WithField("error", err).Errorf("unable to remove the buttons from the slack message")
		return
	}

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
		return
	}
}
