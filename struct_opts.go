package cli

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
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
// The restriction to [BuiltinOpType] allows us to guarantee type safety
// despite linking the resulting Opts to the memory addresses of the struct
// via reflection and unsafe pointer casting.
//
// Calling [Opt.WithValue] on a resulting opt, for example as is done in
// the default Parse implementation, will actually update the corresponding
// struct field directly.
func StructOpts(s any) ([]*Opt, error) {
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
		uPtr := fVal.Addr().UnsafePointer()
		opt, err := cliTagOpt(f.Tag.Get("cli"), uPtr)
		if err != nil {
			return nil, err
		}
		if opt != nil {
			opts = append(opts, opt)
		}
	}
	return opts, nil
}

func cliTagOpt(tag string, p unsafe.Pointer) (*Opt, error) {
	if tag == "" {
		return nil, nil
	}
	tag = strings.TrimSpace(tag)
	opt := &Opt{
		Link: p,
	}

	n := len(tag)
	i := 0
	hasType := false
	for i < n {
		key, rest, ok := strings.Cut(tag[i:], "=")
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
			switch rest {
			case "bool":
				opt.Type = Bool
			case "int":
				opt.Type = Int
			case "float":
				opt.Type = Float
			case "string":
				opt.Type = String
			default:
				return nil, fmt.Errorf("%w: unsupported type: %q", ErrTagParseError, rest)
			}
			hasType = true

		case "desc":
			opt.Description = rest
		case "default":
			if !hasType {
				return nil, fmt.Errorf("%w: default must come after type", ErrTagParseError, key)
			}
			switch opt.Type {
			case Bool:
				v, err := strconv.ParseBool(rest)
				if err != nil {
					return nil, fmt.Errorf("%w: invalid bool %q", ErrTagParseError, rest)
				}
				var a any = v
				opt.Default = &a
			case Int:
				v, err := strconv.ParseInt(rest, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("%w: invalid int %q", ErrTagParseError, rest)
				}
				var a any = int(v)
				opt.Default = &a
			case Float:
				var f float64
				if _, err := fmt.Sscanf(rest, "%f", &f); err != nil {
					return nil, fmt.Errorf("%w: invalid float %q: %w", ErrTagParseError, rest, err)
				}
				var a any = f
				opt.Default = &a
			case String:
				var a any = rest
				opt.Default = &a
			default:
				return nil, fmt.Errorf("%w: unsupported type: %q", ErrTagParseError, rest)
			}
			// TODO
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
	for j < len(v) {
		c := v[j]
		j++
		switch c {
		case '\'':
			if !escaped {
				return j, v[1 : j-1], nil
			}
			escaped = false
		case '\\':
			escaped = !escaped
		default:
			escaped = false
		}

	}
	return 0, "", ErrTagParseError
}
