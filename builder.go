package cli

func NewCommand(name string) *Command {
	cmd := &Command{}
	return NewCommandAt(&cmd, name)
}

func NewCommandAt(p **Command, name string) *Command {
	tmp := &Command{Name: name}
	*p = tmp
	return *p
}

func (cmd *Command) WithSynopsis(s string) *Command {
	cmd.Synopsis = s
	return cmd
}

func (cmd *Command) WithDescription(d string) *Command {
	cmd.Description = d
	return cmd
}

func (cmd *Command) WithAliases(als ...string) *Command {
	cmd.Aliases = append(cmd.Aliases, als...)
	return cmd
}

func (cmd *Command) WithSuppressedOpts(opts ...string) *Command {
	if cmd.InvalidOpts == nil {
		cmd.InvalidOpts = map[string]bool{}
	}
	for _, opt := range opts {
		cmd.InvalidOpts[opt] = true
	}
	return cmd
}

func (cmd *Command) WithSubs(subs ...*Command) *Command {
	for _, sub := range subs {
		sub.Parent = cmd
		cmd.Children = append(cmd.Children, sub)
	}
	return cmd
}

func (cmd *Command) WithOpts(opts ...*Opt) *Command {
	for _, opt := range opts {
		opt.Parent = cmd
		cmd.Opts = append(cmd.Opts, opt)
	}
	return cmd
}

func (cmd *Command) WithParse(pf ParseFunc) *Command {
	cmd.Hooks.Parse = pf
	return cmd
}

func (cmd *Command) WithRun(rf RunFunc) *Command {
	cmd.Hooks.Run = rf
	return cmd
}

func (cmd *Command) WithExit(f ExitFunc) *Command {
	cmd.Hooks.Exit = f
	return cmd
}

func (cmd *Command) WithUsage(f UsageFunc) *Command {
	cmd.Hooks.Usage = f
	return cmd
}
