package ch3

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"exercises/util"
)

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
	GiB // 1073741824
	TiB // 1099511627776             (exceeds 1 << 32)
	PiB // 1125899906842624
	EiB // 1152921504606846976
	//ZiB // 1180591620717411303424    (exceeds 1 << 64)
	//YiB // 1208925819614629174706176
)

// Exercises 练习
type Exercises struct {
}

// Task1
// 1. 如果f函数返回的是无限制的float64值，那么SVG文件可能输出无效的多边形元素（虽然许多SVG渲染器会妥善处理这类问题）。修改程序跳过无效的多边形。
// 2. 试验math包中其他函数的渲染图形。你是否能输出一个egg box、moguls或a saddle图案?
// 3. 根据高度给每个多边形上色，那样峰值部将是红色（#ff0000），谷部将是蓝色（#0000ff）。
// 4. 参考1.7节Lissajous例子的函数，构造一个web服务器，用于计算函数曲面然后返回SVG数据给客户端。服务器必须设置Content-Type头部。
// go run main.go -ch 3 -task 1 zs eggBox
// go run main.go -http=8080 后请求 http://localhost:8080/ch3/surface?shape=eggBox
func (e *Exercises) Task1() {
	file, err := util.NewFile(1, filepath.Join("file", os.Args[5]+".svg"))
	if err != nil {
		return
	}
	defer file.Close()
	if len(os.Args) > 6 {
		Surface(file, os.Args[6])
	} else {
		Surface(file, "")
	}

}

// Task2 to Task1
func (e *Exercises) Task2() {
	e.Task1()
}

// Task3 to Task1
func (e *Exercises) Task3() {
	e.Task1()
}

// Task4 to Task1
func (e *Exercises) Task4() {
	e.Task1()
}

// Task5
// 1. 实现一个彩色的Mandelbrot图像，使用image.NewRGBA创建图像，使用color.RGBA或color.YCbCr生成颜色。
// 2. 升采样技术可以降低每个像素对计算颜色值和平均值的影响。简单的方法是将每个像素分成四个子像素，实现它。
// 3. 另一个生成分形图像的方式是使用牛顿法来求解一个复数方程，例如$z^4-1=0$。每个起点到四个根的迭代次数对应阴影的灰度。方程根对应的点用颜色表示。
// 4. 通过提高精度来生成更多级别的分形。使用四种不同精度类型的数字实现相同的分形：complex64、complex128、big.Float和big.Rat。（后面两种类型在math/big包声明。Float是有指定限精度的浮点数；Rat是无限精度的有理数。）它们间的性能和内存使用对比如何？当渲染图可见时缩放的级别是多少？
// 5. 编写一个web服务器，用于给客户端生成分形的图像。运行客户端通过HTTP参数指定x、y和zoom参数。
// go run main.go -ch 3 -task 5 sc
// go run main.go -http=8080 后请求 http://localhost:8080/ch3/mandelbort
func (e *Exercises) Task5() {
	file, err := util.NewFile(1, filepath.Join("file", os.Args[5]+".png"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	Mandelbrot(file)
}

// Task6 to Task5
func (e *Exercises) Task6() {
	e.Task5()
}

// Task7 to Task5
func (e *Exercises) Task7() {
	e.Task5()
}

// Task8 to Task5
func (e *Exercises) Task8() {
	e.Task5()
}

// Task9 to Task5
func (e *Exercises) Task9() {
	e.Task5()
}

// Task10
// 1. 编写一个非递归版本的comma函数，使用bytes.Buffer代替字符串链接操作。
// 2. 完善comma函数，以支持浮点数处理和一个可选的正负号的处理。
// go run main.go -ch 3 -task 10 1231231231232
func (e *Exercises) Task10() {
	var buf bytes.Buffer
	str := os.Args[5]
	start := 0

	// 正负号
	switch str[start] {
	case '+':
		start++
	case '-':
		start++
		buf.WriteByte('-')
	}

	// 小数点在首位
	switch str[start] {
	case '.':
		start++
		buf.WriteByte('0')
	}

	// 获取小数点的位置
	end := strings.Index(str, ".")

	// 整数部分
	for i := start; i < end; i++ {
		buf.WriteByte(str[i])
		if end-i != 1 && (end-1-i)%3 == 0 {
			buf.WriteByte(',')
		}
	}

	// 小数部分
	if end != -1 {
		buf.WriteString(str[end:])
	}
	fmt.Println(buf.String())
}

// Task11 to Task10
func (e *Exercises) Task11() {
	e.Task10()
}

// Task12 编写一个函数，判断两个字符串是否是相互打乱的，也就是说它们有着相同的字符，但是对应不同的顺序。
// go run main.go -ch 3 -task 12 112 121
func (e *Exercises) Task12() {
	mapA := make(map[rune]int)
	for _, s := range os.Args[5] {
		mapA[s]++
	}
	for _, s := range os.Args[6] {
		mapA[s]--
	}
	for _, val := range mapA {
		if val != 0 {
			fmt.Println(false)
			return
		}
	}
	fmt.Println(true)
}

// Task13 编写KB、MB的常量声明，然后扩展到YB。
func (e *Exercises) Task13() {
	fmt.Println(KiB)
	fmt.Println(MiB)
	fmt.Println(GiB)
	fmt.Println(TiB)
	fmt.Println(PiB)
	fmt.Println(EiB)
	//fmt.Println(ZiB)
	//fmt.Println(YiB)
}
