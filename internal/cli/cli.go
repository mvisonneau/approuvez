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
	app.Usage = "Command line helper to obtain live confirmation from people in a blocking fashion"
	app.EnableBashCompletion = true

	app.Flags = cli.FlagsByName{
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
		&cli.BoolFlag{
			Name:    "tls-disable",
			EnvVars: []string{"APPROUVEZ_TLS_DISABLE"},
			Usage:   "disable mutual tls for gRPC transmissions (use with care!)",
		},
		&cli.StringFlag{
			Name:    "tls-ca-cert",
			EnvVars: []string{"APPROUVEZ_TLS_CA_CERT"},
			Usage:   "TLS CA certificate `path`",
		},
		&cli.StringFlag{
			Name:    "tls-cert",
			EnvVars: []string{"APPROUVEZ_TLS_CERT"},
			Usage:   "TLS certificate `path`",
		},
		&cli.StringFlag{
			Name:    "tls-key",
			EnvVars: []string{"APPROUVEZ_TLS_KEY"},
			Usage:   "TLS key `path`",
		},
	}

	app.Commands = cli.CommandsByName{
		{
			Name:   "ask",
			Usage:  "send a message to someone and wait for a response",
			Action: cmd.ExecWrapper(cmd.Ask),
			Flags: cli.FlagsByName{
				&cli.StringFlag{
					Name:    "endpoint",
					Aliases: []string{"e"},
					EnvVars: []string{"APPROUVEZ_SERVER_ENDPOINT"},
					Usage:   "server `endpoint` to connect upon",
					Value:   "127.0.0.1:8443",
				},
				&cli.StringFlag{
					Name:    "user",
					Aliases: []string{"u"},
					EnvVars: []string{"APPROUVEZ_USER"},
					Usage:   "`email or slack ID` of a person that should review the message",
				},
				&cli.StringFlag{
					Name:    "message",
					Aliases: []string{"m"},
					EnvVars: []string{"APPROUVEZ_MESSAGE"},
					Usage:   "`message` to display on Slack",
				},
				&cli.StringFlag{
					Name:    "link-name",
					EnvVars: []string{"APPROUVEZ_LINK_NAME"},
					Usage:   "`name` of a link button to append to the message",
				},
				&cli.StringFlag{
					Name:    "link-url",
					EnvVars: []string{"APPROUVEZ_LINK_URL"},
					Usage:   "`url` of a link button to append to the message",
				},
			},
		},
		{
			Name:   "serve",
			Usage:  "run the server thing",
			Action: cmd.ExecWrapper(cmd.Serve),
			Flags: cli.FlagsByName{
				&cli.StringFlag{
					Name:    "slack-token",
					EnvVars: []string{"APPROUVEZ_SLACK_TOKEN"},
					Usage:   "`token` to use in order to authenticate requests against slack",
				},
				&cli.StringFlag{
					Name:    "listen-address",
					EnvVars: []string{"APPROUVEZ_LISTEN_ADDRESS"},
					Usage:   "`token` to use in order to authenticate requests against slack",
					Value:   ":8443",
				},
			},
		},
	}

	app.Metadata = map[string]interface{}{
		"startTime": start,
	}

	return
}
