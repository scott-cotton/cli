package cli

import (
	"errors"
	"fmt"
)

func (cmd *Command) Exit(cc *Context, err error) int {
	if cmd.Hooks.Exit != nil {
		return cmd.Hooks.Exit(cc, err)
	}
	if err == nil {
		return 0
	}
	if !errors.Is(err, ErrUsage) {
		fmt.Fprintf(cc.Err, "%v\n", err.Error())
	}
	var xc ExitCodeErr
	if errors.As(err, &xc) {
		return int(xc)
	}
	return 1
}
