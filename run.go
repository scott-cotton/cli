package cli

import (
	"errors"
	"fmt"
	"os"
)

func (cmd *Command) Run(cc *Context, args []string) error {
	if cmd.Hooks.Run != nil {
		return cmd.Hooks.Run(cc, args)
	}
	if len(cmd.Children) == 0 {
		return ErrNoCommandProvided
	}
	args, err := cmd.Parse(cc, args)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return ErrNoCommandProvided
	}
	sub := cmd.FindSub(cc, args[0])
	if sub == nil {
		return fmt.Errorf("%w: %q", ErrNoSuchCommand, args[0])
	}
	err = sub.Run(cc, args[1:])
	if errors.Is(err, ErrUsage) {
		sub.Usage(cc, err)
	}
	os.Exit(sub.Exit(cc, err))
	return err
}

func (cmd *Command) FindSub(cc *Context, sub string) *Command {
	for _, c := range cmd.Children {
		if c.Name == sub {
			return c
		}
		for _, al := range c.Aliases {
			if al != sub {
				continue
			}
			return c
		}
	}
	return nil
}
