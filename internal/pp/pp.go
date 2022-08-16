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

type CustomPrinter func(val any) string

func Custom(printer func(val any) string) CustomPrinter {
	return CustomPrinter(printer)
}

type context struct {
	depth  int
	custom []CustomPrinter
}

var (
	stringerTyper = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	timeType      = reflect.TypeOf(time.Time{})
)

func (ctx context) print(v reflect.Value) string {
	if v.CanInterface() {
		i := v.Interface()
		for _, p := range ctx.custom {
			out := p(i)
			if out != "" {
				return ansi.Cyan.String(out)
			}
		}
	}
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
		return ansi.Blue.String(fmt.Sprintf("%v", v)) + fmt.Sprintf(" (%s)", v.Type())
	}
	return ctx.sprintf("{%s %s} %v", v.Type(), k, v)
}

func (ctx context) sprintf(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func Print(v ...any) {
	ctx := context{
		depth:  1,
		custom: []CustomPrinter{customPrinter},
	}
	vals := []any{}
	for _, one := range v {
		if tp, ok := one.(CustomPrinter); ok {
			ctx.custom = append(ctx.custom, tp)
		} else {
			vals = append(vals, one)
		}
	}
	for _, val := range vals {
		fmt.Println(ctx.print(reflect.ValueOf(val)))
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
		sb.WriteString(ansi.Black.String("0x"))
		sb.WriteString(ansi.Green.String(fmt.Sprintf("%x", out)))
		sb.Write(ansi.ResetBytes)
		sb.Write(ansi.BlackBytes)
		if len(b) > max {
			sb.WriteString("...")
		}
		sb.WriteString(" (")
		sb.WriteString(sizestr.Bytes(len(b)).String())
		sb.WriteString(")")
		sb.Write(ansi.ResetBytes)
		return sb.String()
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

var C = ansi.Attribute("test")

func (ctx context) highlightString(s string) string {
	sb := strings.Builder{}
	for _, r := range s {
		if r == '\n' {
			sb.WriteString(ansi.White.String("⏎") + "\n")
		} else if r == '\t' {
			sb.WriteString(ansi.White.String("⇥"))
		} else if r > 32 && r <= unicode.MaxASCII {
			sb.WriteString(ansi.Green.String(string(r)))
		} else {
			sb.WriteString(ansi.White.String("·"))
		}
	}
	return sb.String()
}

func (ctx context) printString(s string) string {
	if s == "" {
		return ""
	}
	h := ctx.highlightString(s)
	if strings.Contains(s, "\n") {
		return "\n" + indent(h)
	}
	return fmt.Sprintf(`"%s"`, h)
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
		outf := indent(fmt.Sprintf("%s: %s", ft.Name, outv))
		fmt.Fprintf(&sb, "%s\n", outf)
	}
	sb.WriteRune('}')
	return sb.String()
}

func isASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
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
