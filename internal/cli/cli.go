package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/mvisonneau/approuvez/internal/cmd"
	"github.com/urfave/cli/v2"
)

// Run handles the instanciation of the CLI application
func Run(version string, args []string) {
	err := NewApp(version, time.Now()).Run(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// NewApp configures the CLI application
func NewApp(version string, start time.Time) (app *cli.App) {
	app = cli.NewApp()
	app.Name = "approuvez"
	app.Version = version
	app.Usage = "Command line helper to obtain live confirmation from relevant people"
	app.EnableBashCompletion = true

	app.Flags = cli.FlagsByName{
		&cli.StringFlag{
			Name:    "endpoint,e",
			EnvVars: []string{"APPROUVEZ_ENDPOINT"},
			Usage:   "websocket `endpoint` to connect in order to listen for events",
		},
		&cli.StringFlag{
			Name:    "token,t",
			EnvVars: []string{"APPROUVEZ_TOKEN"},
			Usage:   "`token` to use in order to authenticate against the endpoint",
		},
		&cli.StringFlag{
			Name:    "slack-token",
			EnvVars: []string{"APPROUVEZ_SLACK_TOKEN"},
			Usage:   "`token` to use in order to authenticate requests against slack",
		},
		&cli.StringFlag{
			Name:    "slack-message",
			EnvVars: []string{"APPROUVEZ_SLACK_MESSAGE"},
			Usage:   "`message` to display on Slack",
		},
		&cli.StringFlag{
			Name:    "slack-channel",
			EnvVars: []string{"APPROUVEZ_SLACK_CHANNEL"},
			Usage:   "slack `channel` to write the message onto",
		},
		&cli.StringFlag{
			Name:    "triggerrer",
			EnvVars: []string{"APPROUVEZ_TRIGGERRER"},
			Usage:   "`email or slack ID` of the person trigerring the message",
		},
		&cli.StringSliceFlag{
			Name:    "reviewer",
			EnvVars: []string{"APPROUVEZ_REVIEWER"},
			Usage:   "`email or slack ID` of a person that should review the message",
		},
		&cli.IntFlag{
			Name:    "required-approvals",
			EnvVars: []string{"APPROUVEZ_REQUIRED_APPROVALS"},
			Usage:   "`amount` of approvals required to consider it approved (default to all defined reviewers)",
		},
		&cli.StringFlag{
			Name:    "log-level",
			EnvVars: []string{"APPROUVEZ_LOG_LEVEL"},
			Usage:   "log `level` (debug,info,warn,fatal,panic)",
			Value:   "info",
		},
		&cli.StringFlag{
			Name:    "log-format",
			EnvVars: []string{"APPROUVEZ_LOG_FORMAT"},
			Usage:   "log `format` (json,text)",
			Value:   "text",
		},
	}

	app.Action = cmd.ExecWrapper(cmd.Run)

	app.Metadata = map[string]interface{}{
		"startTime": start,
	}

	return
}
