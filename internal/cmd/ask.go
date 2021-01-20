package cmd

import (
	"context"
	"strings"
	"time"

	"github.com/mvisonneau/approuvez/pkg/certs"
	"github.com/mvisonneau/approuvez/pkg/client"
	pb "github.com/mvisonneau/approuvez/pkg/protobuf"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

// Ask sends a message to someone and blocks until a response is provided
func Ask(ctx *cli.Context) (int, error) {
	cfg := configureClient(ctx)

	conn := Dial(cfg)
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

// Dial ..
func Dial(cfg client.Config) *grpc.ClientConn {
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

	log.WithField("endpoint", cfg.Endpoint).Debug("initializing gRPC connection to the server..")
	conn, err := grpc.Dial(
		cfg.Endpoint,
		tlsDialOption,
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		log.WithField("endpoint", cfg.Endpoint).WithField("error", err).Fatalf("could not connect to the server")
	}

	log.Debug("gRPC connection established..")
	return conn
}
