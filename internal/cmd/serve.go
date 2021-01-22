package cmd

import (
	"net"
	"net/http"

	"github.com/mvisonneau/approuvez/pkg/certs"
	pb "github.com/mvisonneau/approuvez/pkg/protobuf"
	"github.com/mvisonneau/approuvez/pkg/server"
	"github.com/soheilhy/cmux"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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
	httpListener := m.Match(cmux.HTTP1Fast())
	grpcListener := m.Match(cmux.Any())

	srv := server.New(cfg.SlackToken)

	g := new(errgroup.Group)
	g.Go(func() error { return httpServe(httpListener, srv) })
	g.Go(func() error { return grpcServe(grpcListener, srv, cfg.TLS) })
	g.Go(func() error { return m.Serve() })

	log.Infof("started multiplexed HTTP/gRPC server on %s", cfg.ListenAddress)

	return 0, g.Wait()
}

func httpServe(l net.Listener, srv *server.Server) error {
	r := mux.NewRouter()
	r.HandleFunc("/callbacks/slack", srv.HandleSlackCallback).Methods("POST")

	// TODO: Implement better healthchecks..
	r.HandleFunc("/health/live", func(rw http.ResponseWriter, r *http.Request) { rw.WriteHeader(http.StatusOK) })
	r.HandleFunc("/health/ready", func(rw http.ResponseWriter, r *http.Request) { rw.WriteHeader(http.StatusOK) })

	s := &http.Server{Handler: r}
	return s.Serve(l)
}

func grpcServe(l net.Listener, srv *server.Server, tlsConfig certs.Config) error {
	var s *grpc.Server
	if tlsConfig.Disable {
		s = grpc.NewServer()
	} else {
		serverCerts, err := certs.LoadServerCertificates(tlsConfig)
		if err != nil {
			log.Fatal(err)
		}
		s = grpc.NewServer(grpc.Creds(serverCerts))
	}

	pb.RegisterApprouvezServer(s, srv)
	return s.Serve(l)
}
