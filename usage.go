package cli

import (
	"errors"
	"fmt"
	"strings"
	"text/tabwriter"
)

func (cmd *Command) Usage(cc *Context, err error) {
	if cmd.Hooks.Usage != nil {
		cmd.Hooks.Usage(cc, err)
		return
	}
	w := cc.Out
	if err != nil {
		w = cc.Err
	}
	fmt.Fprintf(w, cmd.Description)
	fmt.Fprintln(w)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "synopsis: %s\n", cmd.Synopsis)
	if len(cmd.Children) != 0 {
		fmt.Fprintf(w, "\n\ncommands:\n")
		tw := tabwriter.NewWriter(w, 1, 4, 2, ' ', 0)
		for _, cmd := range cmd.Children {
			fmt.Fprintf(tw, "\t")
			if cmd.Synopsis != "" {
				name, _, ok := strings.Cut(cmd.Synopsis, " ")
				if ok {
					fmt.Fprintf(tw, "\t%s\t%s", name, strings.TrimSpace(cmd.Synopsis))
				} else {
					fmt.Fprintf(tw, "\t%q\t", cmd.Name)
				}
			}
			fmt.Fprintln(tw)
		}
		tw.Flush()
	}
	path := cmd.Path()
	for _, c := range path {
		available := "available "
		if c == cmd {
			available = ""
		}
		name := c.Name
		if c.Parent == nil && len(path) == 1 {
			name = ""
		}
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintf(w, "%s%s options:", available, name)
		if len(c.Opts) == 0 {
			fmt.Fprintf(w, " (none)\n")
		}
		tw := tabwriter.NewWriter(w, 0, 1, 1, ' ', tabwriter.TabIndent)
		for _, co := range c.Opts {
			fmt.Fprintf(tw, "\n\t%s\t\t%s", co.FormatFlag(), co.FormatDesc())
		}
		fmt.Fprintln(w)
		tw.Flush()
	}

	if errors.Is(err, ErrUsage) {
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, err.Error())
	}

}

func (o *Opt) FormatFlag() string {
	b := &strings.Builder{}
	b.WriteByte('-')
	b.WriteString(o.Name)
	for _, al := range o.Aliases {
		b.WriteString(", -" + al)
	}
	return b.String()
}

func (o *Opt) FormatDesc() string {
	return o.Description + "\t" + o.Type.String()
}
