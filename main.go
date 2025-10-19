package cli

import (
	"context"
	"errors"
	"os"
)

// Main runs cmd
func Main(cmd *Command) {
	MainContext(context.Background(), cmd)
}

// MainContext runs cmd with a given go context.
func MainContext(ctx context.Context, cmd *Command) {
	cc := &Context{
		Out: os.Stdout,
		Err: os.Stderr,
		In:  os.Stdin,
		Go:  ctx,
		Env: os.Environ(),
	}
	err := cmd.Run(cc, os.Args[1:])
	if errors.Is(err, ErrUsage) {
		cmd.Usage(cc, err)
	}
	os.Exit(cmd.Exit(cc, err))
}
