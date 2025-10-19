package cli

import (
	"context"
	"os"
)

// DefaultContext returns a default cli Context, using
// context.Background, stdin, stdout, stderr, and [os.Environ()].
func DefaultContext() *Context {
	return &Context{
		Out: os.Stdout,
		Err: os.Stderr,
		In:  os.Stdin,
		Env: os.Environ(),
		Go:  context.Background(),
	}
}
