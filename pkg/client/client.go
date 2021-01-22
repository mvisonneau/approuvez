package client

import (
	"context"
	"strings"
	"time"

	"github.com/mvisonneau/approuvez/pkg/certs"
	pb "github.com/mvisonneau/approuvez/pkg/protobuf"
	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
)

// Config ..
type Config struct {
	Endpoint string
	User     string
	Message  string
	LinkName string
	LinkURL  string
	TLS      certs.Config
}

// Client ..
type Client struct {
	RPC pb.ApprouvezClient
}

// Dial ..
func Dial(cfg Config) *grpc.ClientConn {
	var tlsDialOption grpc.DialOption
	if cfg.TLS.Disable {
		tlsDialOption = grpc.WithInsecure()
	} else {
		clientCerts, err := certs.LoadClientCertificates(
			cfg.TLS,
			cfg.Endpoint[:strings.IndexByte(cfg.Endpoint, ':')],
		)
		if err != nil {
			log.Fatal(err)
		}
		tlsDialOption = grpc.WithTransportCredentials(clientCerts)
	}

	log.WithField("endpoint", cfg.Endpoint).Debug("establishing gRPC connection to the server..")
	conn, err := grpc.Dial(
		cfg.Endpoint,
		tlsDialOption,
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		log.WithField("endpoint", cfg.Endpoint).WithField("error", err).Fatalf("could not connect to the server")
	}

	log.Debug("gRPC connection established")
	return conn
}

// ListenForSlackResponses ..
func (c *Client) ListenForSlackResponses(req *pb.SlackUserRequest) (int, error) {
	stream, err := c.RPC.CreateStream(context.Background(), req)
	if err != nil {
		log.WithField("error", err).Errorf("connection failed")
		return 1, nil
	}

	log.Infof("message sent, waiting for user's descision")

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
