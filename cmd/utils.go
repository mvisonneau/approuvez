package cmd

import (
	"os"
	"time"

	"github.com/mvisonneau/approuvez/lib/client"
	"github.com/mvisonneau/go-helpers/logger"
	"github.com/urfave/cli"

	log "github.com/sirupsen/logrus"
)

var start time.Time

func configure(ctx *cli.Context) (c *client.Client, err error) {
	start = ctx.App.Metadata["startTime"].(time.Time)

	lc := &logger.Config{
		Level:  ctx.GlobalString("log-level"),
		Format: ctx.GlobalString("log-format"),
	}

	if err = lc.Configure(); err != nil {
		return
	}

	for _, i := range []string{"endpoint", "slack-token", "slack-message", "slack-channel", "triggerrer"} {
		assertStringVariableDefined(ctx, i, ctx.GlobalString(i))
	}
	assertStringSliceVariableNotEmpty(ctx, "reviewer", ctx.GlobalStringSlice("reviewer"))

	return client.NewClient(&client.NewClientInput{
		SlackChannel:      ctx.GlobalString("slack-channel"),
		SlackMessage:      ctx.GlobalString("slack-message"),
		SlackToken:        ctx.GlobalString("slack-token"),
		WebsocketEndpoint: ctx.GlobalString("endpoint"),
		Triggerrer:        ctx.GlobalString("triggerrer"),
		Reviewers:         ctx.GlobalStringSlice("reviewer"),
		RequiredApprovals: ctx.GlobalInt("required-approvals"),
	})
}

func exit(exitCode int, err error) *cli.ExitError {
	defer log.WithFields(
		"execution-duration": time.Since(start),
	).Debug("exiting..")

	if err != nil {
		log.Error(err.Error())
	}

	return cli.NewExitError("", exitCode)
}

// ExecWrapper gracefully logs and exits our `run` functions
func ExecWrapper(f func(ctx *cli.Context) (int, error)) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		return exit(f(ctx))
	}
}

func assertStringVariableDefined(ctx *cli.Context, k, v string) {
	if len(v) == 0 {
		cli.ShowAppHelp(ctx)
		log.Errorf("'--%s' must be set!", k)
		os.Exit(2)
	}
}

func assertStringSliceVariableNotEmpty(ctx *cli.Context, k string, v []string) {
	if len(v) == 0 {
		cli.ShowAppHelp(ctx)
		log.Errorf("'--%s' must be set at least once!", k)
		os.Exit(2)
	}
}
