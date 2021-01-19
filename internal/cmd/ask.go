package cmd

import (
	"context"
	"time"

	"github.com/mvisonneau/approuvez/pkg/client"
	pb "github.com/mvisonneau/approuvez/pkg/protobuf/approuvez"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

// Ask sends a message to someone and blocks until a response is provided
func Ask(ctx *cli.Context) (int, error) {
	cfg := configureClient(ctx)

	log.WithField("endpoint", cfg.Endpoint).Debug("initializing gRPC connection to the server..")
	conn, err := grpc.Dial(cfg.Endpoint, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.WithField("endpoint", cfg.Endpoint).WithField("error", err).Fatalf("could not connect to the server")
	}

	log.Debug("gRPC connection established..")
	defer conn.Close()

	c := client.Client{pb.NewApprouvezClient(conn)}
	session, err := c.NewSession(context.TODO(), &pb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Connected with session ID '%s'", session.Id)

	req := &pb.SlackUserRequest{
		Session: session,
		User:    cfg.Reviewer,
		Message: cfg.Message,
	}

	return c.ListenForSlackResponses(req)
}
