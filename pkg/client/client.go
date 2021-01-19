package client

import (
	"context"

	pb "github.com/mvisonneau/approuvez/pkg/protobuf/approuvez"

	log "github.com/sirupsen/logrus"
)

// Config ..
type Config struct {
	Endpoint string
	Message  string
	Reviewer string
}

// Client ..
type Client struct {
	pb.ApprouvezClient
}

// ListenForSlackResponses ..
func (c *Client) ListenForSlackResponses(req *pb.SlackUserRequest) (int, error) {
	stream, err := c.CreateStream(context.Background(), req)
	if err != nil {
		log.WithField("error", err).Errorf("connection failed")
		return 1, nil
	}

	msg, err := stream.Recv()
	if err != nil {
		log.WithField("error", err).Errorf("reading message from server")
		return 1, nil
	}

	log.WithFields(log.Fields{
		"user_id":   msg.User.Id,
		"user_name": msg.User.Name,
		"decision":  msg.Decision,
	}).Infof("received response")

	if msg.Decision == pb.SlackUserResponse_APPROVE {
		return 0, nil
	}

	return 1, nil
}
