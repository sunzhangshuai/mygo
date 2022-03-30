package ch13

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"unsafe"

	"exercises/ch13/bzip"
	"exercises/util"
)

// Exercises 练习
type Exercises struct {
}

// Task1 定义一个深比较函数，对于十亿以内的数字比较，忽略类型差异。
// go run main.go -ch=13 -task=1
func (e *Exercises) Task1() {
	var a int64 = 16
	var b uint = 16
	fmt.Println(Equal(a, b))
}

// Task2 编写一个函数，报告其参数是否为循环数据结构。
// go run main.go -ch=13 -task=2
func (e *Exercises) Task2() {
	type data struct {
		Key   string
		Value string
		Next  *data
	}
	key1 := data{"key1", "key2", nil}
	key2 := data{"key3", "key3", &key1}
	key1.Next = &key2
	fmt.Println(isLoop(reflect.ValueOf(key1)))
}

// Task3 使用sync.Mutex以保证bzip2.writer在多个goroutines中被并发调用是安全的。
// go run main.go -ch=13 -task=3
func (e *Exercises) Task3() {
	file, err := util.NewFile(1, filepath.Join("file", "zs.zip"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	ret, err := http.Get("http://gopl.io/")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	str := util.ReadString(ret.Body)
	w := bzip.NewWriter(file)
	n, err := w.Write([]byte(str))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(n)
}

// Task4 因为C库依赖的限制。 使用os/exec包启动/bin/bzip2命令作为一个子进程，提供一个纯Go的bzip.NewWriter的替代实现
// （译注：虽然是纯Go实现，但是运行时将依赖/bin/bzip2命令，其他操作系统可能无法运行）。
// go run main.go -ch=13 -task=4
func (e *Exercises) Task4() {
	filename := util.GetFileName(1, filepath.Join("file", "zs.html"))
	cmd := exec.Command("bzip2", "-k", filename)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

var ptrStack = make(map[unsafe.Pointer]bool)

type comparison struct {
	x, y unsafe.Pointer
	reflect.Type
}

func Equal(x, y interface{}) bool {
	seen := make(map[comparison]bool)
	return equal(reflect.ValueOf(x), reflect.ValueOf(y), seen)
}

func equal(x, y reflect.Value, seen map[comparison]bool) bool {
	if !x.IsValid() || !y.IsValid() {
		return x.IsValid() == y.IsValid()
	}
	if x.Type() != y.Type() && !strings.Contains(x.Type().String(), "int") || !strings.Contains(y.Type().String(), "int") {
		return false
	}

	// ...cycle check omitted (shown later)...
	if x.CanAddr() && y.CanAddr() {
		xptr := unsafe.Pointer(x.UnsafeAddr())
		yptr := unsafe.Pointer(y.UnsafeAddr())
		if xptr == yptr {
			return true // identical references
		}
		c := comparison{xptr, yptr, x.Type()}
		if seen[c] {
			return true // already seen
		}
		seen[c] = true
	}

	switch x.Kind() {
	case reflect.Bool:
		return x.Bool() == y.Bool()
	case reflect.String:
		return x.String() == y.String()

	// ...numeric cases omitted for brevity...
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var xValue, yValue int64
		if strings.Contains(x.Type().String(), "uint") {
			xValue = int64(x.Uint())
		} else {
			xValue = x.Int()
		}
		if strings.Contains(y.Type().String(), "uint") {
			yValue = int64(y.Uint())
		} else {
			yValue = y.Int()
		}
		return xValue == yValue

	case reflect.Chan, reflect.UnsafePointer, reflect.Func:
		return x.Pointer() == y.Pointer()
	case reflect.Ptr, reflect.Interface:
		return equal(x.Elem(), y.Elem(), seen)
	case reflect.Array, reflect.Slice:
		if x.Len() != y.Len() {
			return false
		}
		for i := 0; i < x.Len(); i++ {
			if !equal(x.Index(i), y.Index(i), seen) {
				return false
			}
		}
		return true

		// ...struct and map cases omitted for brevity...
	}
	panic("unreachable")
}

// isLoop 判断是否为循环数据结构。
func isLoop(data reflect.Value) bool {
	if data.Kind() == reflect.Struct {
		for i := 0; i < data.NumField(); i++ {
			if isLoop(data.Field(i)) {
				return true
			}
		}
	} else if data.Kind() == reflect.Ptr {
		if data.CanAddr() {
			if ptrStack[unsafe.Pointer(data.UnsafeAddr())] {
				return true
			}
			ptrStack[unsafe.Pointer(data.UnsafeAddr())] = true
			defer func() {
				ptrStack[unsafe.Pointer(data.UnsafeAddr())] = false
			}()
		}
		return isLoop(data.Elem())
	}
	return false
}
