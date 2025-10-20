// Package cli command.
//
// In short, this source code
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
//		Debug bool `cli:"name=debug type=bool desc='turn on debugging'"`
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
//		Name  string `cli:"name=n type=string desc=name default=sam"`
//		Level int    `cli:"name=level type=int desc='A\\'s level'"`
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
//		# ./example
//		example is a demo use of github.com/scott-cotton/cli
//
//		synopsis: example run me and see
//
//		commands:
//		    a  a the a command exits code equal to the number of args
//		    b  b cool and use cli
//
//		 options:
//
//		 -debug  turn on debugging bool
//		usage error: no command provided
//
//
//		# ./example -debug
//		example is a demo use of github.com/scott-cotton/cli
//
//		synopsis: example run me and see
//
//		commands:
//		    a  a the a command exits code equal to the number of args
//		    b  b cool and use cli
//
//		 options:
//
//		 -debug  turn on debugging bool
//		usage error: no command provided
//
//
//		# ./example -no-debug
//		example is a demo use of github.com/scott-cotton/cli
//
//		synopsis: example run me and see
//
//		commands:
//		    a  a the a command exits code equal to the number of args
//		    b  b cool and use cli
//
//		 options:
//
//		 -debug  turn on debugging bool
//		usage error: no command provided
//
//
//		# ./example -nodebug
//		example is a demo use of github.com/scott-cotton/cli
//
//		synopsis: example run me and see
//
//		commands:
//		    a  a the a command exits code equal to the number of args
//		    b  b cool and use cli
//
//		 options:
//
//		 -debug  turn on debugging bool
//		usage error: unknown option: "nodebug"
//
//
//		# ./example c
//		example is a demo use of github.com/scott-cotton/cli
//
//		synopsis: example run me and see
//
//		commands:
//		    a  a the a command exits code equal to the number of args
//		    b  b cool and use cli
//
//		 options:
//
//		 -debug  turn on debugging bool
//		usage error: no such command: "c"
//
//
//		# ./example a -h
//		synopsis: a the a command exits code equal to the number of args
//
//		available example options:
//
//		 -debug  turn on debugging bool
//
//		a options:
//
//		 -n      name      (default sam) string
//		 -level  A's level int
//		usage error: unknown option: "h"
//
//
//		# ./example a -debug
//		should exit 0
//
//
//		# ./example a x -debug
//		should exit 1
//
//
//		# ./example a x -n wilma -level 2 y
//		should exit 2
//
//
//		# ./example b
//		b is a subcommand
//
//		synopsis: b cool and use cli
//
//		available example options:
//
//		 -debug  turn on debugging bool
//
//		b options:
//
//		 -e  -e key=val sets key to val <func>
//		usage error: please supply some -e flags or args
//
//
//		# ./example b -e x=y
//		args: []
//		env:
//		        x: y
//
//
//		# ./example b arg0 -e x
//		b is a subcommand
//
//		synopsis: b cool and use cli
//
//		available example options:
//
//		 -debug  turn on debugging bool
//
//		b options:
//
//		 -e  -e key=val sets key to val <func>
//		usage error: -e expected key=value
//
//
//		# ./example b arg0 -e x=y -e x2=y2
//		args: [arg0]
//		env:
//		        x: y
//		        x2: y2
//		# ./example b arg0 -e x=y arg1 -debug
//	     args: [arg0 arg1]
//	     env:
//	       x: y
package main
