package ch7

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"exercises/util"
)

// Exercises 练习
type Exercises struct {
}

// Task1 使用来自ByteCounter的思路，实现一个针对单词和行数的计数器。你会发现bufio.ScanWords非常的有用。
// go run main.go -ch=7 -task=1
func (e *Exercises) Task1() {
	file, err := util.GetFile(1, filepath.Join("file", "text.html"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	var counter MyCounter
	byt := make([]byte, 1<<20)
	nr, err := file.Read(byt)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	counter.Write(byt[:nr])
	fmt.Println(counter)
}

// Task2 写一个带有如下函数签名的函数CountingWriter，传入一个io.Writer接口类型，返回一个把原来的Writer封装在里面的新的Writer类型和一个表示新的写入字节数的int64类型指针。
// go run main.go -ch=7 -task=2
func (e *Exercises) Task2() {
	rw, iLen := CountingWriter(os.Stdout)
	_, err := rw.Write([]byte("hello,world\n"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(*iLen)
}

// Task3 为在gopl.io/ch4/treesort（§4.4）中的*tree类型实现一个String方法去展示tree类型的值序列。
// go run main.go -ch=7 -task=3
func (e *Exercises) Task3() {
	values := make([]int, 0, 100)
	for i := 0; i < 100; i++ {
		values = append(values, i)
	}
	tree1 := Sort(values)
	fmt.Println(tree1.String())
}

// Task4 strings.NewReader函数通过读取一个string参数返回一个满足io.Reader接口类型的值（和其它值）。
// 实现一个简单版本的NewReader，用它来构造一个接收字符串输入的HTML解析器（§5.2）
// go run main.go -ch=7 -task=4
func (e *Exercises) Task4() {
	file, err := util.GetFile(1, filepath.Join("file", "text.html"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	fmt.Println(HtmlParse(util.ReadString(file)))
}

// Task5 io包里面的LimitReader函数接收一个io.Reader接口类型的r和字节数n，
// 并且返回另一个从r中读取字节但是当读完n个字节后就表示读到文件结束的Reader。实现这个LimitReader函数
// go run main.go -ch=7 -task=5
func (e *Exercises) Task5() {
	s := "dsfsdfsdghtgyhyjuyjffdgdf"
	r := strings.NewReader(s)
	nr := LimitReader(r, 2)
	byts := make([]byte, 2)
	n, err := nr.Read(byts)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(byts[:n]))
}

// Task6 对tempFlag加入支持开尔文温度。
// go run main.go -ch=7 -task=6 -celsius=18C
func (e *Exercises) Task6() {
	fmt.Println(fg)
}

// Task7 解释为什么帮助信息在它的默认值是20.0没有包含°C的情况下输出了°C。
// go run main.go -ch=7 -task=7
func (e *Exercises) Task7() {
	fmt.Println("tempconv.Celsius 的 String 方法会加上 °C")
}

// Task8 很多图形界面提供了一个有状态的多重排序表格插件：
// 主要的排序键是最近一次点击过列头的列，第二个排序键是第二最近点击过列头的列，等等。
// 定义一个sort.Interface的实现用在这样的表格中。比较这个实现方式和重复使用sort.Stable来排序的方式。
func (e *Exercises) Task8() {
}

// Task9 使用html/template包（§4.6）替代printTracks将tracks展示成一个HTML表格。
// 将这个解决方案用在前一个练习中，让每次点击一个列的头部产生一个HTTP请求来排序这个表格。
func (e *Exercises) Task9() {
}

// Task10 sort.Interface类型也可以适用在其它地方。
// 编写一个IsPalindrome(s sort.Interface) bool函数表明序列s是否是回文序列，
// 换句话说反向排序不会改变这个序列。假设如果!s.Less(i, j) && !s.Less(j, i)则索引i和j上的元素相等。
func (e *Exercises) Task10() {
}
