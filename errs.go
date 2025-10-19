package cli

import (
	"errors"
	"fmt"
)

var (
	ErrUsage             = errors.New("usage error")
	ErrUnimplemented     = errors.New("unimplemented")
	ErrNoCommandProvided = fmt.Errorf("%w: no command provided", ErrUsage)
	ErrNoSuchCommand     = fmt.Errorf("%w: no such command", ErrUsage)
	ErrUnknownOption     = fmt.Errorf("%w: unknown option", ErrUsage)
	ErrOptRequiresValue  = fmt.Errorf("%w: option requires a value", ErrUsage)

	ErrTagParseError = errors.New("tag parse error")
)

type ExitCodeErr int

func (c ExitCodeErr) Error() string {
	return fmt.Sprintf("exit %d", c)
}
