package ch12

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Display(name string, x interface{}, depth int) {
	fmt.Printf("Display %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x), depth)
}

func display(path string, v reflect.Value, depth int) {
	if depth == 0 {
		fmt.Printf("%s = Out of range\n", path)
		return
	}
	depth--
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i), depth)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(fieldPath, v.Field(i), depth)
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path,
				formatAtom(key)), v.MapIndex(key), depth)
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem(), depth)
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display(path+".value", v.Elem(), depth)
		}
	default: // basic types, channels, funcs
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}

// formatAtom formats a value without inspecting its internal structure.
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	case reflect.Interface:
		return formatAtom(v.Elem())
	case reflect.Array:
		strs := make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			strs[i] = formatAtom(v.Index(i))
		}
		return fmt.Sprintf("%s[%s]", v.Type().String(), strings.Join(strs, ", "))
	case reflect.Struct:
		strs := make([]string, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			strs[i] = v.Type().Field(i).Name + ":" + formatAtom(v.Field(i))
		}
		return fmt.Sprintf("%s[%s]", v.Type().String(), strings.Join(strs, ", "))
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}

// Movie 电影
type Movie struct {
	Inter    interface{}       `sexpr:"in"`
	Title    string            `sexpr:"t"`
	Subtitle string            `sexpr:"s"`
	Year     int               `sexpr:"y"`
	Color    bool              `sexpr:"c"`
	Actor    map[string]string `sexpr:"a"`
	Oscars   []string
	Sequel   *string
}
