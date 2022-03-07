package main

import (
	"exercises/ch1"
	"exercises/ch2"
	"exercises/ch3"
	"exercises/ch4"
	"exercises/ch8"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type ChList struct {
	Ch1 *ch1.Exercises
	Ch2 *ch2.Exercises
	Ch3 *ch3.Exercises
	Ch4 *ch4.Exercises
	//Ch5 *ch5.Exercises
	Ch8 *ch8.Exercises
}

// 参数
var chIndex = flag.Int("ch", 1, "章节")
var taskIndex = flag.Int("task", 1, "作业")

// main 练习的入口函数，必传参数 -ch、-task
func main() {
	// 格式化参数
	flag.Parse()
	// 执行作业
	run()
	fmt.Println(os.Args[0])
}

// run 执行作业
func run() {
	reflect.ValueOf(ChList{}).
		FieldByName("Ch" + strconv.Itoa(*chIndex)).
		MethodByName("Task" + strconv.Itoa(*taskIndex)).
		Call([]reflect.Value{})
}
