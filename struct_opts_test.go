package cli

import "testing"

type C struct {
	X bool `cli:"name=x aliases=X desc='hello j'"`
}

func TestConfig(t *testing.T) {
	c := &C{}
	opts, err := StructOpts(c)
	if err != nil {
		t.Error(err)
		return
	}
	cmd := NewCommand("test").WithOpts(opts...)
	_, err = cmd.Parse(DefaultContext(), []string{"-X"})
	if !c.X {
		t.Error("didn't set x\n")
	}
}
