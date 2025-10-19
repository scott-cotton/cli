package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/scott-cotton/cli"
)

func main() {
	cli.MainContext(context.Background(), MainCommand())
}

type MainConfig struct {
	Debug bool `cli:"name=debug type=bool desc='turn on debugging'"`
}

func MainCommand() *cli.Command {
	cfg := &MainConfig{}
	opts, err := cli.StructOpts(cfg)
	if err != nil {
		// yes it should panic for development
		// but nothing would reasonably be shipped with
		// this possibility
		panic(err)
	}
	return cli.NewCommand("example").
		WithSynopsis("example run me and see").
		WithDescription("example is a demo use of github.com/scott-cotton/cli").
		WithOpts(opts...).
		WithSubs(
			ACommand(),
			BCommand(),
		)
}

// The a command
type AConfig struct {
	*cli.Command
	Name  string `cli:"name=n type=string desc=name"`
	Level int    `cli:"name=level type=int desc='A\'s level'"`
}

func ACommand() *cli.Command {
	cfg := &AConfig{}
	opts, err := cli.StructOpts(cfg)
	if err != nil {
		// yes it should panic for development
		// but nothing would reasonably be shipped with
		// this possibility
		panic(err)
	}
	return cli.NewCommandAt(&cfg.Command, "a").
		WithSynopsis("a the a command exits code equal to the number of args").
		WithOpts(opts...).
		WithRun(cfg.run)
}

func (cfg *AConfig) run(cc *cli.Context, args []string) error {
	args, err := cfg.Parse(cc, args)
	if err != nil {
		return err
	}
	fmt.Fprintf(cc.Out, "should exit %d\n", len(args))
	return cli.ExitCodeErr(len(args))
}

// The b command
type BConfig struct {
	*cli.Command
	env map[string]any
}

func (bCfg *BConfig) parseEnv(cc *cli.Context, a string) (any, error) {
	key, val, ok := strings.Cut(a, "=")
	if !ok {
		return nil, fmt.Errorf("%w: -e expected key=value", cli.ErrUsage)
	}
	bCfg.env[key] = val
	return nil, nil
}
func BCommand() *cli.Command {
	cfg := &BConfig{
		env: map[string]any{},
	}
	return cli.NewCommandAt(&cfg.Command, "b").
		WithAliases("bb", "bbb").
		WithSynopsis("b cool and use cli").
		WithDescription("b is a subcommand").
		WithOpts(&cli.Opt{
			Name:        "e",
			Description: "set environment",
			Type:        cli.FuncOpt(cfg.parseEnv),
		}).WithRun(cfg.run)
}

func (b *BConfig) run(cc *cli.Context, args []string) error {
	args, err := b.Parse(cc, args)
	if err != nil {
		fmt.Printf("b parse err %s\n", err.Error())
		return err
	}
	fmt.Fprintf(cc.Out, "args: %v\n", args)
	for k, v := range b.env {
		fmt.Fprintf(cc.Out, "\t%s: %v\n", k, v)
	}
	if len(b.env) == 0 && len(args) == 0 {
		return fmt.Errorf("%w: please supply some -e flags or args", cli.ErrUsage)
	}
	return nil
}
