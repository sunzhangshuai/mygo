// Package decode provides a means for converting Go objects to and
// from S-expressions.
package decode

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"text/scanner"
)

// Unmarshal parses S-expression data and populates the variable
// whose address is in the non-nil pointer out.
func Unmarshal(data []byte, out interface{}) (err error) {
	lex := &Decoder{scan: scanner.Scanner{Mode: scanner.GoTokens}}
	lex.scan.Init(bytes.NewReader(data))
	return lex.marshal(out)
}

// NewDecoder 获取解码器
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decoder 解码器
type Decoder struct {
	r     io.Reader
	scan  scanner.Scanner
	token rune // the current token
}

func (lex *Decoder) Decode(out interface{}) (err error) {
	lex.scan.Init(lex.r)
	return lex.marshal(out)
}

func (lex *Decoder) marshal(out interface{}) (err error) {
	lex.next() // get the first token
	defer func() {
		// NOTE: this is not an example of ideal error handling.
		if x := recover(); x != nil {
			err = fmt.Errorf("error at %s: %v", lex.scan.Position, x)
		}
	}()
	read(lex, reflect.ValueOf(out).Elem())
	return nil
}

func (lex *Decoder) next()        { lex.token = lex.scan.Scan() }
func (lex *Decoder) text() string { return lex.scan.TokenText() }

func (lex *Decoder) consume(want rune) {
	if lex.token != want { // NOTE: Not an example of good error handling.
		panic(fmt.Sprintf("got %q, want %q", lex.text(), want))
	}
	lex.next()
}

// The read function is a decoder for a small subset of well-formed
// S-expressions.  For brevity of our example, it takes many dubious
// shortcuts.
//
// The parser assumes
// - that the S-expression input is well-formed; it does no error checking.
// - that the S-expression input corresponds to the type of the variable.
// - that all numbers in the input are non-negative decimal integers.
// - that all keys in ((key value) ...) struct syntax are unquoted symbols.
// - that the input does not contain dotted lists such as (1 2 . 3).
// - that the input does not contain Lisp reader macros such 'x and #'x.
//
// The reflection logic assumes
// - that v is always a variable of the appropriate type for the
//   S-expression value.  For example, v must not be a boolean,
//   interface, channel, or function, and if v is an array, the input
//   must have the correct number of elements.
// - that v in the top-level call to read has the zero value of its
//   type and doesn't need clearing.
// - that if v is a numeric variable, it is a signed integer.
func read(lex *Decoder, v reflect.Value) {
	switch lex.token {
	case scanner.Ident:
		// The only valid identifiers are
		// "nil" and struct field names.
		if lex.text() == "nil" {
			v.Set(reflect.Zero(v.Type()))
			lex.next()
			return
		}
		if lex.text() == "t" {
			v.SetBool(true)
			lex.next()
			return
		}
	case scanner.String:
		s, _ := strconv.Unquote(lex.text()) // NOTE: ignoring errors
		if v.Type().Kind() == reflect.Interface {
			lex.next()
			v1 := makeInterface(s)
			read(lex, v1.Elem())
			v.Set(v1.Elem())
			return
		}
		v.SetString(s)
		lex.next()
		return
	case scanner.Int:
		i, _ := strconv.Atoi(lex.text()) // NOTE: ignoring errors
		v.SetInt(int64(i))
		lex.next()
		return
	case scanner.Float:
		i, _ := strconv.ParseFloat(lex.text(), 64) // NOTE: ignoring errors
		v.SetFloat(i)
		lex.next()
		return
	case '(':
		lex.next()
		readList(lex, v)
		lex.next() // consume ')'
		return
	}

	panic(fmt.Sprintf("unexpected token %q", lex.text()))
}

func makeInterface(token string) reflect.Value {
	switch token {
	case "int":
		var a int
		return reflect.New(reflect.TypeOf(a))
	}
	panic(fmt.Sprintf("unexpected token %q", token))
}

// readList 读取列表
func readList(lex *Decoder, v reflect.Value) {
	switch v.Kind() {
	case reflect.Array: // (item ...)
		for i := 0; !endList(lex); i++ {
			read(lex, v.Index(i))
		}

	case reflect.Slice: // (item ...)
		for !endList(lex) {
			item := reflect.New(v.Type().Elem()).Elem()
			read(lex, item)
			v.Set(reflect.Append(v, item))
		}

	case reflect.Struct: // ((name value) ...)
		nameAlias := make(map[string]string)
		for i := 0; i < v.NumField(); i++ {
			name := v.Type().Field(i).Name
			tag := v.Type().Field(i).Name
			if v.Type().Field(i).Tag.Get("sexpr") != "" {
				tag = v.Type().Field(i).Tag.Get("sexpr")
			}
			nameAlias[tag] = name
		}

		for !endList(lex) {
			lex.consume('(')
			if lex.token != scanner.Ident {
				panic(fmt.Sprintf("got token %q, want field name", lex.text()))
			}
			name := nameAlias[lex.text()]
			lex.next()
			read(lex, v.FieldByName(name))
			lex.consume(')')
		}

	case reflect.Map: // ((key value) ...)
		v.Set(reflect.MakeMap(v.Type()))
		for !endList(lex) {
			lex.consume('(')
			key := reflect.New(v.Type().Key()).Elem()
			read(lex, key)
			value := reflect.New(v.Type().Elem()).Elem()
			read(lex, value)
			v.SetMapIndex(key, value)
			lex.consume(')')
		}

	default:
		panic(fmt.Sprintf("cannot decode list into %v", v.Type()))
	}
}

// endList 判断list是否结束
func endList(lex *Decoder) bool {
	switch lex.token {
	case scanner.EOF:
		panic("end of file")
	case ')':
		return true
	}
	return false
}

//!-readlist
