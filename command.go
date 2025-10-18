package cli

func (cmd *Command) Root() *Command {
	if cmd.Parent == nil {
		return cmd
	}
	return cmd.Parent.Root()
}

func (cmd *Command) OptMap() map[string]*Opt {
	res := map[string]*Opt{}
	cmd.PutMap(res)
	return res
}

func (cmd *Command) PutMap(m map[string]*Opt) map[string]*Opt {
	for i := range cmd.Opts {
		o := cmd.Opts[i]
		m[o.Name] = o
		for _, al := range o.Aliases {
			m[al] = o
		}
	}
	return m
}

func (cmd *Command) Path() []*Command {
	if cmd.Parent == nil {
		return []*Command{cmd}
	}
	return append(cmd.Parent.Path(), cmd)
}

func (cmd *Command) PutOptsAll(dst map[string]*Opt) map[string]*Opt {
	res := map[string]*Opt{}
	for _, o := range cmd.Path() {
		res = o.PutMap(res)
		for k := range o.InvalidOpts {
			io := res[k]
			if io == nil {
				continue
			}
			for _, al := range io.Aliases {
				delete(res, al)
			}
			delete(res, k)
		}
	}
	return res
}

func (cmd *Command) AllOpts() map[string]*Opt {
	return cmd.PutOptsAll(map[string]*Opt{})
}
