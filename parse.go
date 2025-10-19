package cli

import (
	"errors"
	"fmt"
	"strings"
)

// Parse parses arguments, parsing any [Opt]s and
// returning all non-option arguments as the arguments
// for cmd.
func (cmd *Command) Parse(cc *Context, args []string) ([]string, error) {
	if cmd.Hooks.Parse != nil {
		return cmd.Hooks.Parse(cc, args)
	}
	if len(cmd.Children) == 0 {
		return cmd.parse(cc, args, true)
	}
	return cmd.parse(cc, args, false)
}

func (cmd *Command) parse(cc *Context, args []string, all bool) ([]string, error) {
	d := cmd.AllOpts()
	res := []string{}
	hasDD := false
	skip := -1
	var errs error
	for i, arg := range args {
		if i == skip {
			skip = -1
			continue
		}
		if !all && len(res) > 0 {
			res = append(res, arg)
			continue
		}
		if hasDD {
			res = append(res, arg)
			continue
		}
		if arg == "--" {
			hasDD = true
			res = append(res, arg)
			continue
		}
		if arg == "" {
			res = append(res, arg)
			continue
		}
		if arg[0] != '-' {
			res = append(res, arg)
			continue
		}
		arg = arg[1:]
		if arg == "" {
			res = append(res, "-")
			continue
		}
		if arg != "" && arg[0] == '-' {
			arg = arg[1:]
		}
		name, rest, ok := strings.Cut(arg, "=")
		if ok {
			opt := d[name]
			if opt == nil {
				errs = errors.Join(errs, fmt.Errorf("%w: %q", ErrUnknownOption, name))
				continue
			}
			v, err := opt.Type.Parse(cc, rest)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			}
			opt.WithValue(v)
		}
		opt := d[arg]
		flip := false
		if opt == nil {
			if strings.HasPrefix(arg, "no-") {
				opt = d[arg[3:]]
				if opt != nil {
					flip = true
				} else {
					errs = errors.Join(errs, fmt.Errorf("%w: %q", ErrUnknownOption, arg))
					continue
				}
			} else {
				errs = errors.Join(errs, fmt.Errorf("%w: %q", ErrUnknownOption, arg))
				continue
			}
		}
		if opt.Type == Bool {
			t := true
			if flip {
				t = !t
			}
			opt.WithValue(t)
			continue
		}
		if opt.Type.ArgRequired() {
			if i == len(args)-1 {
				errs = errors.Join(errs, fmt.Errorf("%w: %s", ErrOptRequiresValue, opt.Name))
				continue
			}
			skip = i + 1
			v, err := opt.Type.Parse(cc, args[skip])
			if err != nil {
				errs = errors.Join(errs, err)
			}
			opt.WithValue(v)
			continue
		}
		var x any = opt
		opt.Value = &x
	}
	return res, errs
}
