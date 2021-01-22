package cmd

import (
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func NewTestContext() (ctx *cli.Context, flags *flag.FlagSet) {
	app := cli.NewApp()
	app.Name = "approuvez"

	app.Metadata = map[string]interface{}{
		"startTime": time.Now(),
	}

	flags = flag.NewFlagSet("test", flag.ContinueOnError)
	ctx = cli.NewContext(app, flags, nil)

	return
}

func TestExit(t *testing.T) {
	err := exit(20, fmt.Errorf("test"))
	assert.Equal(t, "", err.Error())
	assert.Equal(t, 20, err.ExitCode())
}

func TestExecWrapper(t *testing.T) {
	ctx, flags := NewTestContext()
	flags.String("log-format", "text", "")
	flags.String("log-level", "debug", "")

	function := func(ctx *cli.Context) (int, error) {
		return 0, nil
	}

	assert.Equal(t, exit(function(ctx)), ExecWrapper(function)(ctx))
}
