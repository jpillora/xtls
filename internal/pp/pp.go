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
			if strings.Contains(out, "\n") {
				out = "\n" + indent(out)
			}
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
		return fmt.Sprintf("%s %s", out, ctx.printType(v.Type()))
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
		s := fmt.Sprintf("%v", v)
		return fmt.Sprintf("%s %s", ansi.Blue.String(s), ctx.printType(v.Type()))
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

func (ctx context) printType(t reflect.Type) string {
	return grey(fmt.Sprintf("(%v)", t))
}

func (ctx context) printByteSlice(v reflect.Value) string {
	// byte slice
	b := v.Interface().([]byte)
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
	sb.Write(ansi.ResetBytes)
	sb.Write(ansi.DimBytes)
	sb.Write(ansi.WhiteBytes)
	sb.WriteString("0x")
	sb.Write(ansi.ResetBytes)
	sb.WriteString(ansi.Green.String(fmt.Sprintf("%x", out)))
	sb.Write(ansi.DimBytes)
	sb.Write(ansi.WhiteBytes)
	if len(b) > max {
		sb.WriteString("...")
	}
	sb.WriteString(" (")
	sb.WriteString(sizestr.Bytes(len(b)).String())
	sb.WriteString(")")
	sb.Write(ansi.ResetBytes)
	return sb.String()
}

func (ctx context) printSlice(s reflect.Value) string {
	t := s.Type()
	if t.Elem().Kind() == reflect.Uint8 {
		return ctx.printByteSlice(s)
	}
	if s.Len() == 0 {
		return ""
	}
	sb := strings.Builder{}
	sb.WriteRune('[')
	sb.WriteRune('\n')
	for i := 0; i < s.Len(); i++ {
		fv := s.Index(i)
		index := grey(fmt.Sprintf("#%d ", i+1))
		out := indent(index + ctx.print(fv))
		fmt.Fprintf(&sb, "%s\n", out)
	}
	sb.WriteRune(']')
	return sb.String()
}

func grey(s string) string {
	return fmt.Sprintf("%s%v%s",
		ansi.Set(ansi.Reset, ansi.Dim, ansi.White),
		s,
		ansi.Set(ansi.Reset),
	)
}

func pink(s string) string {
	return fmt.Sprintf("%s%v%s",
		ansi.Set(ansi.Reset, ansi.Bright, ansi.Red),
		s,
		ansi.Set(ansi.Reset),
	)
}

func (ctx context) highlightString(s string) string {
	sb := strings.Builder{}
	for _, r := range s {
		if r == '\n' {
			sb.WriteString(grey("⏎") + "\n")
		} else if r == '\t' {
			sb.WriteString(grey("⇥"))
		} else if r == ' ' {
			sb.WriteString(ansi.Green.String("·"))
		} else if r > 32 && r <= unicode.MaxASCII {
			sb.WriteString(ansi.Green.String(string(r)))
		} else {
			sb.WriteString(grey("·"))
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
		outf := indent(fmt.Sprintf("%s: %s", pink(ft.Name), outv))
		fmt.Fprintf(&sb, "%s\n", outf)
	}
	sb.WriteRune('}')
	return sb.String()
}

func isASCII(s string) bool {
	for _, r := range s {
		if r < 0 || r > unicode.MaxASCII {
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
