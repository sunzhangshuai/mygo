package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"exercises/ch1"
	"exercises/ch10"
	"exercises/ch12"
	"exercises/ch13"
	"exercises/ch2"
	"exercises/ch3"
	"exercises/ch4"
	"exercises/ch5"
	"exercises/ch6"
	"exercises/ch7"
	"exercises/ch7/eval"
	"exercises/ch8"
	"exercises/ch9"
)

type ChList struct {
	Ch1  *ch1.Exercises
	Ch2  *ch2.Exercises
	Ch3  *ch3.Exercises
	Ch4  *ch4.Exercises
	Ch5  *ch5.Exercises
	Ch6  *ch6.Exercises
	Ch7  *ch7.Exercises
	Ch8  *ch8.Exercises
	Ch9  *ch9.Exercises
	Ch10 *ch10.Exercises
	Ch12 *ch12.Exercises
	Ch13 *ch13.Exercises
}

// 参数
var chIndex = flag.Int("ch", 1, "章节")
var taskIndex = flag.Int("task", 1, "作业")
var httpFlag = flag.String("http", "", "http服务端口")
var celsiusFlag = ch7.CelsiusFlag("celsius", 20.0, "温度")

// main 练习的入口函数，必传参数 -ch、-task
func main() {
	// 格式化参数
	flag.Parse()

	// 有http就启动http
	if *httpFlag != "" {
		startHttp(*httpFlag)
		return
	}
	// 执行作业
	run()
	fmt.Println(os.Args[0])
}

// run 执行作业
func run() {
	ch7.SetCelsius(celsiusFlag)
	reflect.ValueOf(ChList{}).
		FieldByName("Ch" + strconv.Itoa(*chIndex)).
		MethodByName("Task" + strconv.Itoa(*taskIndex)).
		Call([]reflect.Value{})
}

// startHttp 启动http服务
func startHttp(port string) {
	http.HandleFunc("/ch1/lissajous", func(writer http.ResponseWriter, request *http.Request) {
		cycles := request.URL.Query().Get("cycles")
		parseFloat, _ := strconv.ParseFloat(cycles, 8)
		writer.Header().Set("Content-Type", "image/gif")
		fmt.Println(ch1.Lissajous(writer, func() uint8 {
			return 1
		}, parseFloat))
	})
	http.HandleFunc("/ch3/surface", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "image/svg+xml")
		ch3.Surface(writer, request.URL.Query().Get("shape"))
	})
	http.HandleFunc("/ch3/mandelbrot", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "image/png")
		ch3.Mandelbrot(writer)
	})
	http.HandleFunc("/ch4/template", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(ch4.RunTemplate(writer))

	})
	http.HandleFunc("/ch7/sort", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(ch7.FmtTracks(writer, request.URL.Query().Get("order")))

	})

	// ch7/database
	http.HandleFunc("/ch7/database/list", ch7.DB.List)
	http.HandleFunc("/ch7/database/add", ch7.DB.Add)
	http.HandleFunc("/ch7/database/delete", ch7.DB.Delete)
	http.HandleFunc("/ch7/database/update", ch7.DB.Update)
	http.HandleFunc("/ch7/database/get", ch7.DB.Get)

	http.HandleFunc("/ch7/eval", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(eval.Run(writer, request.URL.Query().Get("expr")))
	})

	http.HandleFunc("/ch12/unpack", func(writer http.ResponseWriter, request *http.Request) {
		var data struct {
			Labels     []string `http:"l" valid:"2"`
			MaxResults int      `http:"max" valid:"1"`
			Exact      bool     `http:"x"`
		}
		if err := ch12.Unpack(request, &data); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest) // 400
			return
		}
		fmt.Fprintf(writer, "Search: %+v\n", data)
	})

	log.Fatal(http.ListenAndServe("localhost:"+port, nil))
}
