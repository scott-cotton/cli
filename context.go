package cli

import (
	"context"
	"os"
)

func DefaultContext() *Context {
	return &Context{
		Out: os.Stdout,
		Err: os.Stderr,
		In:  os.Stdin,
		Env: os.Environ(),
		Go:  context.Background(),
	}
}
