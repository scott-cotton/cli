package cli

import "fmt"

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
		return fmt.Errorf("%w: %s expects a subcommand", CommandUsageErr(cmd), cmd.Name)
	}
	sub := cmd.FindSub(cc, args[0])
	if sub == nil {
		return fmt.Errorf("%w: %s does not have command %q", CommandUsageErr(cmd), cmd.Name, args[0])
	}
	return sub.Run(cc, args[1:])
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
