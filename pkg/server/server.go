package server

import (
	"context"
	"fmt"

	"github.com/mvisonneau/approuvez/pkg/certs"
	pb "github.com/mvisonneau/approuvez/pkg/protobuf"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"google.golang.org/grpc/peer"
)

// Config ..
type Config struct {
	SlackToken    string
	ListenAddress string
	TLS           certs.Config
}

// Server handles necessary components to run the server side of the app
type Server struct {
	pb.UnimplementedApprouvezServer

	Slack    *slack.Client
	Sessions Sessions
}

// Session ..
type Session struct {
	ID     string
	Stream pb.Approuvez_CreateStreamServer
	Error  chan error
}

// Sessions ..
type Sessions map[string]*Session

// New ..
func New(slackToken string) *Server {
	return &Server{
		Slack:    slack.New(slackToken),
		Sessions: make(Sessions),
	}
}

// NewSession ..
func (s *Server) NewSession(ctx context.Context, _ *pb.Empty) (*pb.Session, error) {
	session := &pb.Session{
		Id: uuid.New().String(),
	}

	p, _ := peer.FromContext(ctx)
	log.WithFields(log.Fields{
		"client_endpoint": p.Addr.String(),
		"session_id":      session.Id,
	}).Info("new session initialized")

	return session, nil
}

// CreateStream ..
func (s *Server) CreateStream(req *pb.SlackUserRequest, stream pb.Approuvez_CreateStreamServer) error {
	sessionID, err := uuid.Parse(req.Session.GetId())
	if err != nil {
		return err
	}

	session := &Session{
		ID:     sessionID.String(),
		Stream: stream,
		Error:  make(chan error),
	}

	s.Sessions[sessionID.String()] = session

	log.WithField("user", req.User).Debug("fetching slack user")
	user, err := s.GetSlackUser(req.User)
	if err != nil {
		return fmt.Errorf("slack error: %v", err)
	}

	msg := SlackMessage{
		SessionID:      sessionID.String(),
		Message:        req.Message,
		LinkButtonName: req.LinkName,
		LinkButtonURL:  req.LinkUrl,
		ActionButtons:  true,
		StatusMessage:  ":grin: waiting for your approval",
		Color:          messageColorDefault,
	}

	if err = s.PromptSlackUser(msg, user.ID); err != nil {
		return err
	}

	return <-session.Error
}
