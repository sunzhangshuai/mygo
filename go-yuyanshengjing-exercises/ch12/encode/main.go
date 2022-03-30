package encode

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// encoder 处理器
type encoder interface {
	encode(buf *bytes.Buffer) error
	encodeIndent(buf *bytes.Buffer, depth int) error
	encodeJson(buf *bytes.Buffer) error
}

// makeEncoder 获取编码器
func makeEncoder(v reflect.Value) (encoder, error) {
	switch v.Kind() {
	case reflect.Invalid:
		return Invalid{v}, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Int{v}, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return Uint{v}, nil

	case reflect.String:
		return String{v}, nil

	case reflect.Ptr:
		return makeEncoder(v.Elem())

	case reflect.Array:
		return Array{v}, nil
	case reflect.Slice: // (value ...)
		if v.IsNil() {
			return Invalid{v}, nil
		}
		return Array{v}, nil

	case reflect.Struct: // ((name value) ...)
		return Struct{v}, nil

	case reflect.Map: // ((key value) ...)
		if v.IsNil() {
			return Invalid{v}, nil
		}
		return Map{v}, nil

	case reflect.Bool:
		return Bool{v}, nil

	case reflect.Float32, reflect.Float64:
		return Float{v}, nil

	case reflect.Complex64, reflect.Complex128:
		return Complex{v}, nil

	case reflect.Interface:
		if v.IsNil() {
			return Invalid{v}, nil
		}
		return Interface{v}, nil

	default: // func
		return nil, fmt.Errorf("unsupported type: %s", v.Type())
	}
}

// isComposite 判断是否为负责数据类型
func isComposite(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Invalid,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.String,
		reflect.Bool,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return false

	case reflect.Ptr, reflect.Interface:
		return isComposite(v.Elem())

	default: // func
		return true
	}
}

// encode 编码
func encode(buf *bytes.Buffer, v reflect.Value) error {
	encoder, err := makeEncoder(v)
	if err != nil {
		return err
	}
	if err = encoder.encode(buf); err != nil {
		return err
	}
	return nil
}

func encodeJson(buf *bytes.Buffer, v reflect.Value) error {
	encoder, err := makeEncoder(v)
	if err != nil {
		return err
	}
	if err = encoder.encodeJson(buf); err != nil {
		return err
	}
	return nil
}

// encodeIndent 格式化编码
func encodeIndent(buf *bytes.Buffer, v reflect.Value, depth int) error {
	encoder, err := makeEncoder(v)
	if err != nil {
		return err
	}
	if err = encoder.encodeIndent(buf, depth); err != nil {
		return err
	}
	return nil
}

func newLine(buf *bytes.Buffer, depth int) {
	buf.WriteByte('\n')
	buf.WriteString(strings.Repeat(" ", depth))
}

// Marshal encodes a Go value in S-expression form.
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalIndent 格式化输出s格式
func MarshalIndent(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encodeIndent(&buf, reflect.ValueOf(v), 0); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalJson 格式化输出s格式
func MarshalJson(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encodeJson(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
