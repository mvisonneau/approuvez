package cli

import (
	"os"
	"time"

	"github.com/mvisonneau/approuvez/cmd"
	"github.com/urfave/cli"
)

// Run handles the instanciation of the CLI application
func Run(version string) {
	NewApp(version, time.Now()).Run(os.Args)
}

// NewApp configures the CLI application
func NewApp(version string, start time.Time) (app *cli.App) {
	app = cli.NewApp()
	app.Name = "approuvez"
	app.Version = version
	app.Usage = "command line helper to obtain live confirmation from relevant people"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "endpoint,e",
			EnvVar: "APPROUVEZ_ENDPOINT",
			Usage:  "websocket `endpoint` to connect in order to listen for events",
		},
		cli.StringFlag{
			Name:   "token,t",
			EnvVar: "APPROUVEZ_TOKEN",
			Usage:  "`token` to use in order to authenticate against the endpoint",
		},
		cli.StringFlag{
			Name:   "slack-token",
			EnvVar: "APPROUVEZ_SLACK_TOKEN",
			Usage:  "`token` to use in order to authenticate requests against slack",
		},
		cli.StringFlag{
			Name:   "slack-message",
			EnvVar: "APPROUVEZ_SLACK_MESSAGE",
			Usage:  "`message` to display on Slack",
		},
		cli.StringFlag{
			Name:   "slack-channel",
			EnvVar: "APPROUVEZ_SLACK_CHANNEL",
			Usage:  "slack `channel` to write the message onto",
		},
		cli.StringFlag{
			Name:   "triggerrer",
			EnvVar: "APPROUVEZ_TRIGGERRER",
			Usage:  "`email or slack ID` of the person trigerring the message",
		},
		cli.StringSliceFlag{
			Name:   "reviewer",
			EnvVar: "APPROUVEZ_REVIEWER",
			Usage:  "`email or slack ID` of a person that should review the message",
		},
		cli.IntFlag{
			Name:   "required-approvals",
			EnvVar: "APPROUVEZ_REQUIRED_APPROVALS",
			Usage:  "`amount` of approvals required to consider it approved (default to all defined reviewers)",
		},
		cli.StringFlag{
			Name:   "log-level",
			EnvVar: "APPROUVEZ_LOG_LEVEL",
			Usage:  "log `level` (debug,info,warn,fatal,panic)",
			Value:  "info",
		},
		cli.StringFlag{
			Name:   "log-format",
			EnvVar: "APPROUVEZ_LOG_FORMAT",
			Usage:  "log `format` (json,text)",
			Value:  "text",
		},
	}

	app.Action = cmd.ExecWrapper(cmd.Run)

	app.Metadata = map[string]interface{}{
		"startTime": start,
	}

	return
}
