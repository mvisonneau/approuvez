package cmd

import (
	"os"
	"time"

	"github.com/mvisonneau/approuvez/pkg/certs"
	"github.com/mvisonneau/approuvez/pkg/client"
	"github.com/mvisonneau/approuvez/pkg/server"
	"github.com/mvisonneau/go-helpers/logger"
	"github.com/urfave/cli/v2"

	log "github.com/sirupsen/logrus"
)

var start time.Time

func configure(ctx *cli.Context) error {
	start = ctx.App.Metadata["startTime"].(time.Time)

	return logger.Configure(logger.Config{
		Level:  ctx.String("log-level"),
		Format: ctx.String("log-format"),
	})
}

func configureClient(ctx *cli.Context) client.Config {
	for _, i := range []string{"endpoint", "message", "user"} {
		assertStringVariableDefined(ctx, i, ctx.String(i))
	}

	return client.Config{
		Endpoint: ctx.String("endpoint"),
		User:     ctx.String("user"),
		Message:  ctx.String("message"),
		LinkName: ctx.String("link-name"),
		LinkURL:  ctx.String("link-url"),
		TLS:      configureTLS(ctx),
	}
}

func configureServer(ctx *cli.Context) server.Config {
	for _, i := range []string{"listen-address", "slack-token"} {
		assertStringVariableDefined(ctx, i, ctx.String(i))
	}

	return server.Config{
		ListenAddress: ctx.String("listen-address"),
		SlackToken:    ctx.String("slack-token"),
		TLS:           configureTLS(ctx),
	}
}

func configureTLS(ctx *cli.Context) certs.Config {
	if !ctx.Bool("tls-disable") {
		for _, i := range []string{"tls-ca-cert", "tls-cert", "tls-key"} {
			assertStringVariableDefined(ctx, i, ctx.String(i))
		}
	}

	return certs.Config{
		Disable: ctx.Bool("tls-disable"),
		CA:      ctx.String("tls-ca-cert"),
		Cert:    ctx.String("tls-cert"),
		Key:     ctx.String("tls-key"),
	}
}

func exit(exitCode int, err error) cli.ExitCoder {
	defer log.WithFields(
		log.Fields{
			"execution-time": time.Since(start),
		},
	).Debug("exited..")

	if err != nil {
		log.Error(err.Error())
	}

	return cli.NewExitError("", exitCode)
}

// ExecWrapper gracefully logs and exits our `run` functions
func ExecWrapper(f func(ctx *cli.Context) (int, error)) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if err := configure(ctx); err != nil {
			return exit(1, err)
		}
		return exit(f(ctx))
	}
}

func assertStringVariableDefined(ctx *cli.Context, k, v string) {
	if len(v) == 0 {
		_ = cli.ShowAppHelp(ctx)
		log.Errorf("'--%s' must be set!", k)
		os.Exit(2)
	}
}

func assertStringSliceVariableNotEmpty(ctx *cli.Context, k string, v []string) {
	if len(v) == 0 {
		_ = cli.ShowAppHelp(ctx)
		log.Errorf("'--%s' must be set at least once!", k)
		os.Exit(2)
	}
}
