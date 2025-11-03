package cli

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// StructOpts uses reflection on structs to create Opts whose values are stored
// in the struct as a field.
//
// Options are created for any field with a `cli:"..."` struct tag.
// The field format permits creating options for [BuiltinOpType] typed
// arguments.  Here is an example:
//
//	type CommandConfig struct {
//	    Debug bool `cli:"name=debug type=bool desc='turn on debugging'"`
//	}
//
// Calling [Opt.WithValue] on a resulting opt, for example as is done in
// the default Parse implementation, will actually update the corresponding
// struct field directly.
func StructOpts(s any) ([]*Opt, error) {
	return StructOptsWithTypes(s, nil)
}

var builtinMap = map[string]OptType{
	"bool":   Bool,
	"string": String,
	"int":    Int,
	"float":  Float,
}

// StructOptsWithTypes
func StructOptsWithTypes(s any, tyMap map[string]OptType) ([]*Opt, error) {
	sMap := map[string]OptType{}
	for k, v := range builtinMap {
		sMap[k] = v
	}
	for k, v := range tyMap {
		sMap[k] = v
	}

	ty := reflect.TypeOf(s)
	if ty == nil {
		return nil, nil
	}
	val := reflect.ValueOf(s)
	switch ty.Kind() {
	case reflect.Struct:
	case reflect.Pointer:
		val = val.Elem()
		ty = ty.Elem()
		if ty.Kind() != reflect.Struct {
			return nil, nil
		}
	default:
		return nil, nil
	}
	n := ty.NumField()
	var opts []*Opt
	for i := range n {
		f := ty.Field(i)
		if f.Anonymous {
			continue
		}
		fVal := val.Field(i)
		if !fVal.CanAddr() {
			continue
		}
		opt, err := cliTagOpt(f.Tag.Get("cli"), fVal, sMap)
		if err != nil {
			return nil, err
		}
		if opt != nil {
			opts = append(opts, opt)
		}
	}
	return opts, nil
}

func cliTagOpt(tag string, fVal reflect.Value, tyMap map[string]OptType) (*Opt, error) {
	if tag == "" {
		return nil, nil
	}
	p := fVal.Addr().UnsafePointer()
	tag = strings.TrimSpace(tag)
	opt := &Opt{
		Link: p,
	}
	hasType := true
	switch fVal.Interface().(type) {
	case bool:
		opt.Type = Bool
	case int:
		opt.Type = Int
	case string:
		opt.Type = String
	case float64:
		opt.Type = Float
	default:
		hasType = false
	}

	n := len(tag)
	i := 0
	for i < n {
		key, _, ok := strings.Cut(tag[i:], "=")
		if !ok {
			return nil, ErrTagParseError
		}
		if i == n-1 {
			return nil, ErrTagParseError
		}
		i += len(key) + 1
		j, rest, err := findRest(tag[i:])
		if err != nil {
			return nil, err
		}
		i += j
		key = strings.TrimSpace(key)
		switch key {
		case "name":
			opt.Name = rest
		case "type":
			if hasType {
				return nil, fmt.Errorf("type specified but inferred (%s)", opt.Type)
			}
			opt.Type = tyMap[rest]
			if opt.Type == nil {
				return nil, fmt.Errorf("%w: unsupported type: %q", ErrTagParseError, rest)
			}
			hasType = true
		case "aliases":
			als := strings.Split(rest, ",")
			for _, al := range als {
				al = strings.TrimSpace(al)
				if al != "" {
					opt.Aliases = append(opt.Aliases, al)
				}
			}
		case "desc":
			opt.Description = rest
		case "default":
			if !hasType {
				return nil, fmt.Errorf("%w: default must come after type for %s", ErrTagParseError, key)
			}
			v, err := opt.Type.Parse(DefaultContext(), rest)
			if err != nil {
				return nil, fmt.Errorf("%w: %w", ErrTagParseError, err)
			}
			opt.Default = &v
			opt = opt.WithValue(v)
		default:
			return nil, fmt.Errorf("%w: unknown tag key %q", ErrTagParseError, key)
		}
	}
	return opt, nil
}

// find value in <key>=<value> where value may be single quoted
// and have backslash escapes of single quotes.
func findRest(v string) (int, string, error) {
	if v == "" {
		return 0, "", nil
	}
	if v[0] != '\'' {
		for j, r := range v {
			if unicode.IsSpace(r) {
				return j, v[0:j], nil
			}
		}
		return len(v), v, nil
	}
	escaped := false
	j := 1
	b := &strings.Builder{}
	for j < len(v) {
		c := v[j]
		j++
		switch c {
		case '\'':
			if !escaped {
				return j, b.String(), nil
			}
			b.WriteByte(c)
			escaped = false
		case '\\':
			if escaped {
				b.WriteByte(c)
			}
			escaped = !escaped
		default:
			b.WriteByte(c)
			escaped = false
		}
	}
	return 0, "", ErrTagParseError
}
