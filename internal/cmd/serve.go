package cmd

import (
	"net"
	"net/http"

	pb "github.com/mvisonneau/approuvez/pkg/protobuf/approuvez"
	"github.com/mvisonneau/approuvez/pkg/server"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// Serve ..
func Serve(ctx *cli.Context) (int, error) {
	cfg := configureServer(ctx)
	listener, err := net.Listen("tcp", cfg.ListenAddress)
	if err != nil {
		log.Fatal(err)
	}

	m := cmux.New(listener)
	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	srv := server.New(cfg.SlackToken)

	g := new(errgroup.Group)
	g.Go(func() error { return grpcServe(grpcListener, srv) })
	g.Go(func() error { return httpServe(httpListener, srv) })
	g.Go(func() error { return m.Serve() })

	log.Infof("started multiplexed HTTP/gRPC server on %s", cfg.ListenAddress)
	g.Wait()

	return 0, nil
}

func httpServe(l net.Listener, srv *server.Server) error {
	r := mux.NewRouter()
	r.HandleFunc("/callbacks/slack", srv.HandleSlackCallback).Methods("POST")

	s := &http.Server{Handler: r}
	return s.Serve(l)
}

func grpcServe(l net.Listener, srv *server.Server) error {
	s := grpc.NewServer()
	pb.RegisterApprouvezServer(s, srv)
	return s.Serve(l)
}
