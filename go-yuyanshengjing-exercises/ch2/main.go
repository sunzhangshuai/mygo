package ch2

import (
	"fmt"
	"strconv"
	"time"

	"exercises/ch2/tempconv"
	"exercises/ch2/unitconv"
	"exercises/util"
)

// Exercises 练习
type Exercises struct {
}

// pc[i] is the population count of i.
var pc = func() (pc [256]byte) {
	for i := range pc {
		// 右移一位的数字已经统计过了，看看最低位是不是1。
		pc[i] = pc[i>>1] + byte(i&1)
	}
	return
}()

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

// Task1 向 tempconv 包添加类型、常量和函数用来处理Kelvin绝对温度的转换，Kelvin 绝对零度是−273.15°C，Kelvin绝对温度1K和摄氏度1°C的单位间隔是一样的。
// go run main.go -ch=2 -task=1
func (e *Exercises) Task1() {
	fmt.Printf("%s == %s\n", tempconv.CToF(tempconv.AbsoluteZeroC), tempconv.AbsoluteZeroF)
	fmt.Printf("%s == %s\n", tempconv.CToK(tempconv.FreezingC), tempconv.FreezingK)
	fmt.Printf("%s == %s\n", tempconv.FToC(tempconv.BoilingF), tempconv.BoilingC)
	fmt.Printf("%s == %s\n", tempconv.FToK(tempconv.FreezingF), tempconv.FreezingK)
	fmt.Printf("%s == %s\n", tempconv.KToC(tempconv.AbsoluteZeroK), tempconv.AbsoluteZeroC)
	fmt.Printf("%s == %s\n", tempconv.KToF(tempconv.BoilingK), tempconv.BoilingF)
}

// Task2 写一个通用的单位转换程序，用类似cf程序的方式从命令行读取参数，如果缺省的话则是从标准输入读取参数，然后做类似Celsius和Fahrenheit的单位转换，长度单位可以对应英尺和米，重量单位可以对应磅和公斤等。
// go run main.go -ch=2 -task=2 100
func (e *Exercises) Task2() {
	var number float64
	var err error
	params := util.GetParam()

	if len(params) == 0 {
		fmt.Println("参数有误，需要数字")
		return
	}
	if number, err = strconv.ParseFloat(params[0], 10); err != nil {
		fmt.Println("参数有误，需要数字")
		return
	}
	fmt.Printf("%s == %s\n", unitconv.Foot(number), unitconv.FtoM(unitconv.Foot(number)))
	fmt.Printf("%s == %s\n", unitconv.Meters(number), unitconv.MtoF(unitconv.Meters(number)))
	fmt.Printf("%s == %s\n", unitconv.Pound(number), unitconv.PtoK(unitconv.Pound(number)))
	fmt.Printf("%s == %s\n", unitconv.Kg(number), unitconv.KtoP(unitconv.Kg(number)))
}

// Task3 重写PopCount函数，用一个循环代替单一的表达式。比较两个版本的性能。
// go run main.go -ch=2 -task=3 999
func (e *Exercises) Task3() {
	maxCount := 100000
	params := util.GetParam()
	if len(params) == 0 {
		fmt.Println("参数有误，需要数字")
		return
	}
	parseUint, err := strconv.ParseUint(params[0], 10, 64)
	if err != nil {
		fmt.Println("参数有误，需要数字")
		return
	}

	start := time.Now()
	for i := 0; i < maxCount; i++ {
		PopCount(parseUint)
	}
	time1 := time.Now()
	fmt.Printf("PopCount cost：%s\n", time1.Sub(start).String())
	for i := 0; i < maxCount; i++ {
		var result byte
		for i := 0; i < 8; i++ {
			result += pc[byte(parseUint>>(i*8))]
		}
	}

	time2 := time.Now()
	fmt.Printf("loop cost：%s\n", time2.Sub(time1).String())
}

// Task4 用移位算法重写PopCount函数，每次测试最右边的1bit，然后统计总数。比较和查表算法的性能差异。
// go run main.go -ch=2 -task=4 999
func (e *Exercises) Task4() {
	maxCount := 100000
	params := util.GetParam()
	if len(params) == 0 {
		fmt.Println("参数有误，需要数字")
		return
	}
	parseUint, err := strconv.ParseUint(params[0], 10, 64)
	if err != nil {
		fmt.Println("参数有误，需要数字")
		return
	}

	start := time.Now()
	for i := 0; i < maxCount; i++ {
		PopCount(parseUint)
	}
	time1 := time.Now()
	fmt.Printf("PopCount cost：%s\n", time1.Sub(start).String())
	for i := 0; i < maxCount; i++ {
		var result byte
		var x = parseUint
		for ; x != 0; x = x >> 1 {
			if x&1 == 1 {
				result++
			}
		}
	}
	time2 := time.Now()
	fmt.Printf("Bit cost：%s\n", time2.Sub(time1).String())
	return
}

// Task5 表达式x&(x-1)用于将x的最低的一个非零的bit位清零。使用这个算法重写PopCount函数，然后比较性能。
// go run main.go -ch=2 -task=5 999
func (e *Exercises) Task5() {
	maxCount := 100000
	params := util.GetParam()
	if len(params) == 0 {
		fmt.Println("参数有误，需要数字")
		return
	}
	parseUint, err := strconv.ParseUint(params[0], 10, 64)
	if err != nil {
		fmt.Println("参数有误，需要数字")
		return
	}

	start := time.Now()
	for i := 0; i < maxCount; i++ {
		PopCount(parseUint)
	}
	time1 := time.Now()
	fmt.Printf("PopCount cost：%s\n", time1.Sub(start).String())
	for i := 0; i < maxCount; i++ {
		var result byte
		var x = parseUint
		for x > 0 {
			x = x & (x - 1)
			result++
		}
	}
	time2 := time.Now()
	fmt.Printf("Bit cost：%s\n", time2.Sub(time1).String())
	return
}
