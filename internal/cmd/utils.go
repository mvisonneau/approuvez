package cmd

import (
	"os"
	"time"

	"github.com/mvisonneau/approuvez/pkg/client"
	"github.com/mvisonneau/go-helpers/logger"
	"github.com/urfave/cli/v2"

	log "github.com/sirupsen/logrus"
)

var start time.Time

func configure(ctx *cli.Context) (c *client.Client, err error) {
	start = ctx.App.Metadata["startTime"].(time.Time)

	// Configure logger
	if err = logger.Configure(logger.Config{
		Level:  ctx.String("log-level"),
		Format: ctx.String("log-format"),
	}); err != nil {
		return
	}

	for _, i := range []string{"endpoint", "slack-token", "slack-message", "slack-channel", "triggerrer"} {
		assertStringVariableDefined(ctx, i, ctx.String(i))
	}
	assertStringSliceVariableNotEmpty(ctx, "reviewer", ctx.StringSlice("reviewer"))

	return client.NewClient(&client.NewClientInput{
		SlackChannel:      ctx.String("slack-channel"),
		SlackMessage:      ctx.String("slack-message"),
		SlackToken:        ctx.String("slack-token"),
		WebsocketEndpoint: ctx.String("endpoint"),
		Triggerrer:        ctx.String("triggerrer"),
		Reviewers:         ctx.StringSlice("reviewer"),
		RequiredApprovals: ctx.Int("required-approvals"),
	})
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
