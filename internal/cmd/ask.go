package cmd

import (
	"context"

	"github.com/mvisonneau/approuvez/pkg/client"
	pb "github.com/mvisonneau/approuvez/pkg/protobuf"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Ask sends a message to someone and blocks until a response is provided
func Ask(ctx *cli.Context) (int, error) {
	cfg := configureClient(ctx)

	conn := client.Dial(cfg)
	defer conn.Close()

	c := client.Client{RPC: pb.NewApprouvezClient(conn)}
	session, err := c.RPC.NewSession(context.TODO(), &pb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	log.WithField("session_id", session.Id).Infof("session initiated successfully")

	req := &pb.SlackUserRequest{
		Session:  session,
		User:     cfg.User,
		Message:  cfg.Message,
		LinkName: cfg.LinkName,
		LinkUrl:  cfg.LinkURL,
	}

	return c.ListenForSlackResponses(req)
}
