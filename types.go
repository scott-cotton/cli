package cli

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"unsafe"
)

// A Command represents a Root CLI command or
// one of its sub-commands, recursively.
type Command struct {
	Name        string
	Aliases     []string
	Description string
	Synopsis    string
	Parent      *Command
	Children    []*Command
	Opts        []*Opt
	InvalidOpts map[string]bool // no aliases

	Hooks CommandHooks
}

// CommandHooks provide hooks to override the default
// implementations provided in this package.
type CommandHooks struct {
	Usage UsageFunc
	Run   RunFunc
	Exit  ExitFunc
	Parse ParseFunc
}

// Context contains CLI context: the input, output,
// error output, environment, and the Go context.
type Context struct {
	In       io.ReadCloser
	Out, Err io.WriteCloser
	Env      []string
	Go       context.Context
}

// Func types for Hooks.
type UsageFunc func(*Context, error)
type RunFunc func(*Context, []string) error
type ExitFunc func(*Context, error) int
type ParseFunc func(*Context, []string) ([]string, error)

// Opt is the type of an option
type Opt struct {
	Name        string
	Aliases     []string
	Description string
	Parent      *Command
	Type        OptType
	Default     *any
	Value       *any
	Link        unsafe.Pointer
}

// WithLink must be called with a pointer to a type
// which corresponds to [o.Type].  Once set up in
// this fashion, [WithValue] will update the pointer.
//
// Provided primarily for [StructOpts].
func (o *Opt) WithLink(p unsafe.Pointer) *Opt {
	o.Link = p
	return o
}

func (o *Opt) WithDefault(v any) *Opt {
	o.Default = &v
	return o
}

func (o *Opt) WithValue(v any) *Opt {
	if o.Link != nil {
		switch o.Type {
		case Bool:
			b := v.(bool)
			linkPtr := (*bool)(o.Link)
			*linkPtr = b
		case String:
			s := v.(string)
			linkPtr := (*string)(o.Link)
			*linkPtr = s
		case Int:
			i := v.(int)
			linkPtr := (*int)(o.Link)
			*linkPtr = i
		case Float:
			f := v.(float64)
			linkPtr := (*float64)(o.Link)
			*linkPtr = f
		default:
			panic("opt link")
		}
	}
	o.Value = &v
	return o
}

// OptType is an interface for the type of a Command Opt
type OptType interface {
	// Parse parses v and returns the result or any error in parsing.  The
	// result is stored in the .Value field of the corresponding options
	// with this option type.
	Parse(cc *Context, v string) (any, error)

	// ArgRequired indicates whether the Option requires an argument.  The
	// only provided OptType which does not require an argument is Bool.
	ArgRequired() bool

	// Stringer
	String() string
}

// BuiltinOptType represents
type BuiltinOptType int

const (
	Bool BuiltinOptType = iota
	Int
	Float
	String
)

func (b BuiltinOptType) ArgRequired() bool {
	switch b {
	case Bool:
		return false
	default:
		return true
	}
}

func (b BuiltinOptType) Parse(_ *Context, v string) (any, error) {
	switch b {
	case Bool:
		return strconv.ParseBool(v)
	case Int:
		i64, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		return int(i64), nil
	case Float:
		var f float64
		_, err := fmt.Sscanf(v, "%f", &f)
		if err != nil {
			return 0.0, err
		}
		return f, nil
	case String:
		return v, nil
	default:
		panic("builtin-type")
	}
}

func (b BuiltinOptType) String() string {
	switch b {
	case Bool:
		return "bool"
	case Int:
		return "int"
	case Float:
		return "float"
	case String:
		return "string"
	default:
		panic("builtin-type")
	}
}

// FuncOpt implements [OptType], providing a
// way to associate any option with a user provided
// function.
type FuncOpt func(*Context, string) (any, error)

func (f FuncOpt) Parse(cc *Context, v string) (any, error) {
	return f(cc, v)
}
func (f FuncOpt) ArgRequired() bool {
	return true
}
func (f FuncOpt) String() string {
	return "<func>"
}
