package encode

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// Invalid invalid
type Invalid struct {
	reflect.Value
}

func (i Invalid) encode(buf *bytes.Buffer) error {
	_, err := buf.WriteString("nil")
	return err
}

func (i Invalid) encodeIndent(buf *bytes.Buffer, depth int) error {
	return i.encode(buf)
}

func (i Invalid) encodeJson(buf *bytes.Buffer) error {
	_, err := buf.WriteString("null")
	return err
}

// Int int
type Int struct {
	reflect.Value
}

func (i Int) encode(buf *bytes.Buffer) error {
	_, err := fmt.Fprintf(buf, "%d", i.Int())
	return err
}

func (i Int) encodeIndent(buf *bytes.Buffer, depth int) error {
	return i.encode(buf)
}

func (i Int) encodeJson(buf *bytes.Buffer) error {
	return i.encode(buf)
}

// Uint uint
type Uint struct {
	reflect.Value
}

func (u Uint) encode(buf *bytes.Buffer) error {
	_, err := fmt.Fprintf(buf, "%d", u.Uint())
	return err
}

func (u Uint) encodeIndent(buf *bytes.Buffer, depth int) error {
	return u.encode(buf)
}

func (u Uint) encodeJson(buf *bytes.Buffer) error {
	return u.encode(buf)
}

// String string
type String struct {
	reflect.Value
}

func (s String) encode(buf *bytes.Buffer) error {
	_, err := fmt.Fprintf(buf, "%q", s.String())
	return err
}

func (s String) encodeIndent(buf *bytes.Buffer, depth int) error {
	return s.encode(buf)
}

func (s String) encodeJson(buf *bytes.Buffer) error {
	_, err := fmt.Fprintf(buf, "\"%s\"", strings.Replace(s.String(), "\"", "\\\"", -1))
	return err
}

// Bool 布尔
type Bool struct {
	reflect.Value
}

func (b Bool) encode(buf *bytes.Buffer) error {
	if b.Bool() {
		buf.WriteString("t")
	} else {
		buf.WriteString("nil")
	}
	return nil
}

func (b Bool) encodeIndent(buf *bytes.Buffer, depth int) error {
	return b.encode(buf)
}

func (b Bool) encodeJson(buf *bytes.Buffer) error {
	var err error
	if b.Bool() {
		_, err = buf.WriteString("true")
	} else {
		_, err = buf.WriteString("false")
	}
	return err
}

// Float 浮点数
type Float struct {
	reflect.Value
}

func (f Float) encode(buf *bytes.Buffer) error {
	_, err := fmt.Fprintf(buf, "%f", f.Float())
	return err
}

func (f Float) encodeIndent(buf *bytes.Buffer, depth int) error {
	return f.encode(buf)
}

func (f Float) encodeJson(buf *bytes.Buffer) error {
	return f.encode(buf)
}

// Complex 复数
type Complex struct {
	reflect.Value
}

func (c Complex) encode(buf *bytes.Buffer) error {
	_, err := fmt.Fprintf(buf, "#C(%f %f)", real(c.Complex()), imag(c.Complex()))
	return err
}

func (c Complex) encodeIndent(buf *bytes.Buffer, depth int) error {
	return c.encode(buf)
}

func (c Complex) encodeJson(buf *bytes.Buffer) error {
	_, err := fmt.Fprintf(buf, "\"%f+%fi\"", real(c.Complex()), imag(c.Complex()))
	return err
}

// Array array
type Array struct {
	reflect.Value
}

func (a Array) encode(buf *bytes.Buffer) error {
	buf.WriteByte('(')
	for i := 0; i < a.Len(); i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		if err := encode(buf, a.Index(i)); err != nil {
			return err
		}
	}
	buf.WriteByte(')')
	return nil
}

func (a Array) encodeIndent(buf *bytes.Buffer, depth int) error {
	buf.WriteByte('(')
	for i := 0; i < a.Len(); i++ {
		newLine(buf, depth+1)
		if err := encodeIndent(buf, a.Index(i), depth+1); err != nil {
			return err
		}
	}
	newLine(buf, depth)
	buf.WriteByte(')')
	return nil
}

func (a Array) encodeJson(buf *bytes.Buffer) error {
	buf.WriteByte('[')
	for i := 0; i < a.Len(); i++ {
		if err := encodeJson(buf, a.Index(i)); err != nil {
			return err
		}
		if i < a.Len()-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(']')
	return nil
}

// Map map
type Map struct {
	reflect.Value
}

func (m Map) encode(buf *bytes.Buffer) error {
	buf.WriteByte('(')
	for i, key := range m.MapKeys() {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteByte('(')
		if err := encode(buf, key); err != nil {
			return err
		}
		buf.WriteByte(' ')
		if err := encode(buf, m.MapIndex(key)); err != nil {
			return err
		}
		buf.WriteByte(')')
	}
	buf.WriteByte(')')
	return nil
}

func (m Map) encodeIndent(buf *bytes.Buffer, depth int) error {
	buf.WriteByte('(')
	for _, key := range m.MapKeys() {
		newLine(buf, depth+1)
		buf.WriteByte('(')
		if err := encodeIndent(buf, key, depth+1); err != nil {
			return err
		}
		buf.WriteByte(' ')
		if err := encodeIndent(buf, m.MapIndex(key), depth+1); err != nil {
			return err
		}
		buf.WriteByte(')')
	}
	newLine(buf, depth)
	buf.WriteByte(')')
	return nil
}

func (m Map) encodeJson(buf *bytes.Buffer) error {
	buf.WriteByte('{')
	for i, key := range m.MapKeys() {
		var nBuf bytes.Buffer
		var err error
		if err = encodeJson(&nBuf, key); err != nil {
			return err
		}
		if isComposite(key) {
			_, err = fmt.Fprintf(buf, "\"%s\"", strings.Replace(nBuf.String(), "\"", "\\\"", -1))
		} else {
			_, err = fmt.Fprintf(buf, "%s", nBuf.String())
		}
		if err != nil {
			return err
		}
		buf.WriteString(":")
		if err = encodeJson(buf, m.MapIndex(key)); err != nil {
			return err
		}
		if i < m.Len()-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte('}')
	return nil
}

// Struct 结构体
type Struct struct {
	reflect.Value
}

func (s Struct) encode(buf *bytes.Buffer) error {
	buf.WriteByte('(')
	for i := 0; i < s.NumField(); i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		name := s.Type().Field(i).Name
		if s.Type().Field(i).Tag.Get("sexpr") != "" {
			name = s.Type().Field(i).Tag.Get("sexpr")
		}
		if _, err := fmt.Fprintf(buf, "(%s ", name); err != nil {
			return err
		}
		if err := encode(buf, s.Field(i)); err != nil {
			return err
		}
		buf.WriteByte(')')
	}
	buf.WriteByte(')')
	return nil
}

func (s Struct) encodeIndent(buf *bytes.Buffer, depth int) error {
	buf.WriteByte('(')
	for i := 0; i < s.NumField(); i++ {
		newLine(buf, depth+1)
		if _, err := fmt.Fprintf(buf, "(%s ", s.Type().Field(i).Name); err != nil {
			return err
		}
		if err := encodeIndent(buf, s.Field(i), depth+1); err != nil {
			return err
		}
		buf.WriteByte(')')
	}
	newLine(buf, depth)
	buf.WriteByte(')')
	return nil
}

func (s Struct) encodeJson(buf *bytes.Buffer) error {
	buf.WriteByte('{')
	for i := 0; i < s.NumField(); i++ {
		if _, err := fmt.Fprintf(buf, "\"%s\":", s.Type().Field(i).Name); err != nil {
			return err
		}
		if err := encodeJson(buf, s.Field(i)); err != nil {
			return err
		}
		if i < s.NumField()-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte('}')
	return nil
}

// Interface 接口
type Interface struct {
	reflect.Value
}

func (i Interface) encode(buf *bytes.Buffer) error {
	_, err := fmt.Fprintf(buf, "\"%s\" ", i.Elem().Type().String())
	if err != nil {
		return err
	}
	return encode(buf, i.Elem())
}

func (i Interface) encodeIndent(buf *bytes.Buffer, depth int) error {
	_, err := fmt.Fprintf(buf, "\"%s\" ", i.Elem().Type().String())
	if err != nil {
		return err
	}
	return encodeIndent(buf, i.Elem(), depth)
}

func (i Interface) encodeJson(buf *bytes.Buffer) error {
	return encodeJson(buf, i.Elem())
}
