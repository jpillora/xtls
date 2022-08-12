package pp

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/jpillora/ansi"
	"github.com/jpillora/sizestr"
)

type context struct {
	depth int
}

var (
	stringerTyper = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	timeType      = reflect.TypeOf(time.Time{})
)

func (ctx context) print(v reflect.Value) string {
	if v.Type() == timeType {
		t := v.Interface().(time.Time)
		return ctx.printString(t.Format(time.RFC3339Nano))
	}
	if v.Type().Implements(stringerTyper) {
		s := v.Interface().(fmt.Stringer)
		out := ctx.printString(s.String())

		return fmt.Sprintf("%s (%s)", out, v.Type())
	}
	if v.Kind() == reflect.String {
		return ctx.printString(v.String())
	}
	if v.Kind() == reflect.Pointer {
		s := v.Elem()
		if s.Kind() == reflect.Struct {
			return ctx.printStruct(s)
		}
	}
	if v.Kind() == reflect.Struct {
		return ctx.printStruct(v)
	}
	if v.Kind() == reflect.Slice {
		return ctx.printSlice(v)
	}
	if v.Kind() == reflect.Bool {
		return ansi.Yellow.String(fmt.Sprintf("%v", v))
	}
	k := strings.ToLower(v.Kind().String())
	if strings.HasPrefix(k, "int") || strings.HasPrefix(k, "uint") || strings.HasPrefix(k, "float") {
		return ansi.Blue.String(fmt.Sprintf("%v", v))
	}
	return ctx.sprintf("{%s %s} %v", v.Type(), k, v)
}

func (ctx context) sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func Print(v ...interface{}) {
	ctx := context{
		depth: 1,
	}
	for _, one := range v {
		fmt.Println(ctx.print(reflect.ValueOf(one)))
	}
}

func (ctx context) printSlice(s reflect.Value) string {
	t := s.Type()
	if t.Elem().Kind() == reflect.Uint8 {
		// byte slice
		b := s.Interface().([]byte)
		// check if ascii text
		s := string(b)
		if isASCII(s) {
			return ctx.printString(s)
		}
		// print as partial hex
		out := b
		const max = 16
		if len(b) > max {
			out = out[0:max]
		}
		sb := strings.Builder{}
		sb.WriteString("0x")
		sb.WriteString(ansi.Green.String(fmt.Sprintf("%x", out)))
		if len(b) > max {
			sb.WriteString("...")
		}
		return fmt.Sprintf("%s (%s)", sb.String(), sizestr.Bytes(len(b)))
	}
	if s.Len() == 0 {
		return ""
	}
	sb := strings.Builder{}
	sb.WriteRune('[')
	sb.WriteRune('\n')
	for i := 0; i < s.Len(); i++ {
		fv := s.Index(i)
		out := indent(ctx.print(fv))
		fmt.Fprintf(&sb, "%s\n", out)
	}
	sb.WriteRune(']')
	return sb.String()
}

func (ctx context) printString(s string) string {
	if s == "" {
		return ""
	}
	if strings.Contains(s, "\n") {
		return "\n" + indent(ansi.Green.String(s))
	}
	return fmt.Sprintf(`"%s"`, ansi.Green.String(s))
}

func (ctx context) printStruct(s reflect.Value) string {
	t := s.Type()
	sb := strings.Builder{}
	if name := s.Type().Name(); name != "" {
		sb.WriteString(name)
	}
	sb.WriteRune('{')
	sb.WriteRune('\n')
	ctx.depth++
	for i := 0; i < s.NumField(); i++ {
		fv := s.Field(i)
		ft := t.Field(i)
		outv := ctx.print(fv)
		if outv == "" {
			continue
		}
		outf := indent(fmt.Sprintf("%s = %s", ft.Name, outv))
		fmt.Fprintf(&sb, "%s\n", outf)
	}
	sb.WriteRune('}')
	return sb.String()
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func indent(s string) string {
	parts := strings.Split(s, "\n")
	for i, p := range parts {
		parts[i] = "  " + p
	}
	return strings.Join(parts, "\n")
}
