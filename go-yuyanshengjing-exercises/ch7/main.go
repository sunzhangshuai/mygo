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
// go run main.go -ch=7 -task=5 -celsius=18
func (e *Exercises) Task6() {
	fmt.Println(fg)
}

// Task7 解释为什么帮助信息在它的默认值是20.0没有包含°C的情况下输出了°C。
func (e *Exercises) Task7() {
	fmt.Println("tempconv.Celsius 的 String 方法会加上 °C")
}
