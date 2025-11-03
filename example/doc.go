// Package cli command.
//
// # In short, this source code
//
//	package main
//
//	import (
//		"context"
//		"fmt"
//		"strings"
//
//		"github.com/scott-cotton/cli"
//	)
//
//	func main() {
//		cli.MainContext(context.Background(), MainCommand())
//	}
//
//	type MainConfig struct {
//		Debug bool `cli:"name=debug desc='turn on debugging'"`
//	}
//
//	func MainCommand() *cli.Command {
//		cfg := &MainConfig{}
//		opts, err := cli.StructOpts(cfg)
//		if err != nil {
//			// yes it should panic for development
//			// but nothing would reasonably be shipped with
//			// this possibility
//			panic(err)
//		}
//		return cli.NewCommand("example").
//			WithSynopsis("example run me and see").
//			WithDescription("example is a demo use of github.com/scott-cotton/cli").
//			WithOpts(opts...).
//			WithSubs(
//				ACommand(),
//				BCommand(),
//			)
//	}
//
//	// The a command
//	type AConfig struct {
//		*cli.Command
//		Name  string `cli:"name=n aliases=name,na desc=name default=sam"`
//		Level int    `cli:"name=level aliases=l desc='A\\'s level'"`
//	}
//
//	func ACommand() *cli.Command {
//		cfg := &AConfig{}
//		opts, err := cli.StructOpts(cfg)
//		if err != nil {
//			// yes it should panic for development
//			// but nothing would reasonably be shipped with
//			// this possibility
//			panic(err)
//		}
//		return cli.NewCommandAt(&cfg.Command, "a").
//			WithSynopsis("a the a command exits code equal to the number of args").
//			WithOpts(opts...).
//			WithRun(cfg.run)
//	}
//
//	func (cfg *AConfig) run(cc *cli.Context, args []string) error {
//		args, err := cfg.Parse(cc, args)
//		if err != nil {
//			return err
//		}
//		fmt.Fprintf(cc.Out, "should exit %d\n", len(args))
//		return cli.ExitCodeErr(len(args))
//	}
//
//	// The b command
//	type BConfig struct {
//		*cli.Command
//		env map[string]any
//
//		// cli.FuncOpt must be bound to field with type any
//		E any `cli:"name=e type=env desc='-e key=val sets key to val'"`
//	}
//
//	func (bCfg *BConfig) parseEnv(cc *cli.Context, a string) (any, error) {
//		key, val, ok := strings.Cut(a, "=")
//		if !ok {
//			return nil, fmt.Errorf("%w: -e expected key=value", cli.ErrUsage)
//		}
//		bCfg.env[key] = val
//		return nil, nil
//	}
//	func BCommand() *cli.Command {
//		cfg := &BConfig{
//			env: map[string]any{},
//		}
//		opts, err := cli.StructOptsWithTypes(cfg, map[string]cli.OptType{
//			"env": cli.FuncOpt(cfg.parseEnv),
//		})
//		if err != nil {
//			panic(err)
//		}
//		return cli.NewCommandAt(&cfg.Command, "b").
//			WithAliases("bb", "bbb").
//			WithSynopsis("b cool and use cli").
//			WithDescription("b is a subcommand").
//			WithOpts(opts...).
//			WithRun(cfg.run)
//	}
//
//	func (b *BConfig) run(cc *cli.Context, args []string) error {
//		args, err := b.Parse(cc, args)
//		if err != nil {
//			return err
//		}
//		if len(b.env) == 0 && len(args) == 0 {
//			return fmt.Errorf("%w: please supply some -e flags or args", cli.ErrUsage)
//		}
//		fmt.Fprintf(cc.Out, "args: %v\nenv:\n", args)
//		for k, v := range b.env {
//			fmt.Fprintf(cc.Out, "\t%s: %v\n", k, v)
//		}
//		return nil
//	}
//
// Produces this CLI
//
//	25-11-03 scott@air example % ./example
//	synopsis: example run me and see
//
//	example is a demo use of github.com/scott-cotton/cli
//
//	commands:
//	    a  a the a command exits code equal to the number of args
//	    b  b cool and use cli
//
//	 options:
//
//	 -debug  turn on debugging bool
//
//	usage error: no command provided
//	25-11-03 scott@air example % ./example -debug
//	synopsis: example run me and see
//
//	example is a demo use of github.com/scott-cotton/cli
//
//	commands:
//	    a  a the a command exits code equal to the number of args
//	    b  b cool and use cli
//
//	 options:
//
//	 -debug  turn on debugging bool
//
//	usage error: no command provided
//	25-11-03 scott@air example % ./example -no-debug
//	synopsis: example run me and see
//
//	example is a demo use of github.com/scott-cotton/cli
//
//	commands:
//	    a  a the a command exits code equal to the number of args
//	    b  b cool and use cli
//
//	 options:
//
//	 -debug  turn on debugging bool
//
//	usage error: no command provided
//	25-11-03 scott@air example % ./example -nodebug
//	synopsis: example run me and see
//
//	example is a demo use of github.com/scott-cotton/cli
//
//	commands:
//	    a  a the a command exits code equal to the number of args
//	    b  b cool and use cli
//
//	 options:
//
//	 -debug  turn on debugging bool
//
//	usage error: unknown option: "nodebug"
//	25-11-03 scott@air example % ./example c
//	synopsis: example run me and see
//
//	example is a demo use of github.com/scott-cotton/cli
//
//	commands:
//	    a  a the a command exits code equal to the number of args
//	    b  b cool and use cli
//
//	 options:
//
//	 -debug  turn on debugging bool
//
//	usage error: no such command: "c"
//	25-11-03 scott@air example % ./example a -h
//	synopsis: a the a command exits code equal to the number of args
//
//	available example options:
//
//	 -debug  turn on debugging bool
//
//	a options:
//
//	 -n, -name, -na  name      (default sam) string
//	 -level, -l      A's level int
//
//	usage error: unknown option: "h"
//	25-11-03 scott@air example % ./example a -debug
//	should exit 0
//	25-11-03 scott@air example % /example a x -debug
//	zsh: no such file or directory: /example
//	25-11-03 scott@air example % ./example a x -debug
//	should exit 1
//	25-11-03 scott@air example % ./example a x -n wilma -level 2 y
//	should exit 2
//	25-11-03 scott@air example % ./example b
//	synopsis: b cool and use cli
//
//	b is a subcommand
//
//	available example options:
//
//	 -debug  turn on debugging bool
//
//	b options:
//
//	 -e  -e key=val sets key to val (func)
//
//	usage error: please supply some -e flags or args
//	25-11-03 scott@air example % ./example b -e x=y
//	args: []
//	env:
//	        x: y
//	25-11-03 scott@air example % ./example b arg0 -e x
//	synopsis: b cool and use cli
//
//	b is a subcommand
//
//	available example options:
//
//	 -debug  turn on debugging bool
//
//	b options:
//
//	 -e  -e key=val sets key to val (func)
//
//	usage error: -e expected key=value
//	25-11-03 scott@air example % ./example b arg0 -e x=y -e x2=y2
//	args: [arg0]
//	env:
//	        x: y
//	        x2: y2
//	25-11-03 scott@air example % ./example b arg0 -e x=y arg1 -debug
//	args: [arg0 arg1]
//	env:
//	        x: y
package main
