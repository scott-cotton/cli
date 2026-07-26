package cli

import (
	"errors"
	"testing"
)

type eqConfig struct {
	Addr  string  `cli:"name=addr desc='an address'"`
	N     int     `cli:"name=n desc='a count'"`
	F     float64 `cli:"name=f desc='a ratio'"`
	Debug bool    `cli:"name=debug desc='a switch'"`
}

func eqCmd(t *testing.T, c *eqConfig) *Command {
	t.Helper()
	opts, err := StructOpts(c)
	if err != nil {
		t.Fatal(err)
	}
	return NewCommand("test").WithOpts(opts...)
}

// -name=value must both APPLY the value and leave no error behind. It used to do the
// first and not the second: the option was looked up by the part before the "=", parsed
// and set, and then the whole "name=value" string fell through to the plain-option lookup,
// where it matched nothing and was reported as an unknown option. The caller saw a usage
// error naming a flag that exists — after the flag had already taken effect.
func TestParseEqualsFormIsNotAnUnknownOption(t *testing.T) {
	c := &eqConfig{}
	cmd := eqCmd(t, c)

	args, err := cmd.Parse(DefaultContext(), []string{"--addr=localhost:9797"})
	if err != nil {
		t.Errorf("Parse(--addr=localhost:9797) error = %v, want nil", err)
	}
	if c.Addr != "localhost:9797" {
		t.Errorf("addr = %q, want %q", c.Addr, "localhost:9797")
	}
	if len(args) != 0 {
		t.Errorf("residual args = %v, want none", args)
	}
}

// Every builtin type goes through the same branch, so the "=" form has to reach each one:
// a bool spelled out, an int, a float, a string.
func TestParseEqualsFormForEachType(t *testing.T) {
	c := &eqConfig{}
	cmd := eqCmd(t, c)

	_, err := cmd.Parse(DefaultContext(),
		[]string{"--addr=host:1", "--n=7", "--f=1.5", "--debug=true"})
	if err != nil {
		t.Fatalf("Parse error = %v, want nil", err)
	}
	if c.Addr != "host:1" || c.N != 7 || c.F != 1.5 || !c.Debug {
		t.Errorf("got %+v, want {Addr:host:1 N:7 F:1.5 Debug:true}", *c)
	}
}

// -debug=false is the only way to say "off" in the "=" form, so it must actually turn the
// switch off rather than being read as "the flag is present, hence true".
func TestParseEqualsFormCanClearABool(t *testing.T) {
	c := &eqConfig{Debug: true}
	cmd := eqCmd(t, c)

	if _, err := cmd.Parse(DefaultContext(), []string{"--debug=false"}); err != nil {
		t.Fatalf("Parse error = %v, want nil", err)
	}
	if c.Debug {
		t.Error("debug = true, want false")
	}
}

// The single dash spelling is the same option, and the "=" form works with either.
func TestParseEqualsFormWithASingleDash(t *testing.T) {
	c := &eqConfig{}
	cmd := eqCmd(t, c)

	if _, err := cmd.Parse(DefaultContext(), []string{"-addr=host:2"}); err != nil {
		t.Fatalf("Parse error = %v, want nil", err)
	}
	if c.Addr != "host:2" {
		t.Errorf("addr = %q, want %q", c.Addr, "host:2")
	}
}

// A genuinely unknown option in the "=" form is still refused, and is named by the option
// part alone — the fix must not make every misspelling parse silently.
func TestParseEqualsFormStillRefusesAnUnknownOption(t *testing.T) {
	c := &eqConfig{}
	cmd := eqCmd(t, c)

	_, err := cmd.Parse(DefaultContext(), []string{"--nope=1"})
	if !errors.Is(err, ErrUnknownOption) {
		t.Errorf("Parse(--nope=1) error = %v, want ErrUnknownOption", err)
	}
}

// A value the option's type cannot parse is an error, not a zero silently assigned.
func TestParseEqualsFormReportsABadValue(t *testing.T) {
	c := &eqConfig{N: 3}
	cmd := eqCmd(t, c)

	if _, err := cmd.Parse(DefaultContext(), []string{"--n=lots"}); err == nil {
		t.Error("Parse(--n=lots) error = nil, want a parse error")
	}
	if c.N != 3 {
		t.Errorf("n = %d, want the original 3 — a rejected value was applied anyway", c.N)
	}
}
