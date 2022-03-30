package ch12

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Unpack 解包数据
func Unpack(req *http.Request, ptr interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	// Build map of fields keyed by effective name.
	fields := make(map[string]reflect.Value)
	checks := make(map[string]reflect.Value)
	v := reflect.ValueOf(ptr).Elem() // the struct variable
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("http")
		checkName := tag.Get("valid")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = v.Field(i)
		if checkName != "" {
			checks[checkName] = v.Field(i)
		}
	}

	// Update struct field for each parameter in the request.
	for name, values := range req.Form {
		f := fields[name]
		if !f.IsValid() {
			continue // ignore unrecognized HTTP parameters
		}
		for _, value := range values {
			if f.Kind() == reflect.Slice {
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.Set(reflect.Append(f, elem))
			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}

	// 校验
	for key, value := range checks {
		if !valid(key, value) {
			return fmt.Errorf("%s check fail", value.Type().String())
		}
	}
	return nil
}

// Pack 打包参数
func Pack(ptr interface{}) (string, error) {
	values := url.Values{}

	v := reflect.ValueOf(ptr).Elem() // the struct variable
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		value := v.Field(i)
		tag := fieldInfo.Tag // a reflect.StructTag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		if value.Kind() == reflect.Slice {
			for j := 0; j < value.Len(); j++ {
				val, err := parse(value.Index(j))
				if err != nil {
					return "", err
				}
				values.Add(name, val)
			}
		} else {
			val, err := parse(value)
			if err != nil {
				return "", err
			}
			values.Set(name, val)
		}

	}
	return values.Encode(), nil
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)

	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}

func parse(v reflect.Value) (string, error) {
	switch v.Kind() {
	case reflect.String:
		return v.String(), nil
	case reflect.Int:
		return strconv.Itoa(int(v.Int())), nil
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), nil
	default:
		return "", fmt.Errorf("unsupported kind %s", v.Type())
	}
}

func valid(rule string, v reflect.Value) bool {
	switch rule {
	case "1":
		return v.Int() > 2
	case "2":
		return v.Len() > 2
	default:
		return false
	}
}
