package ch7

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"exercises/ch7/eval"
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

// Task8
// 1. 很多图形界面提供了一个有状态的多重排序表格插件： 主要的排序键是最近一次点击过列头的列，第二个排序键是第二最近点击过列头的列，等等。 定义一个sort.Interface的实现用在这样的表格中。比较这个实现方式和重复使用sort.Stable来排序的方式。
// 2. 使用html/template包（§4.6）替代printTracks将tracks展示成一个HTML表格。将这个解决方案用在前一个练习中，让每次点击一个列的头部产生一个HTTP请求来排序这个表格。
// go run main.go -ch=7 -task=8 Title|Artist|Album|Year|Length
// or go run main.go -http=8080 后请求 http://localhost:8080/ch7/sort
func (e *Exercises) Task8() {
	params := util.GetParam()
	if len(params) == 0 {
		fmt.Println("请输入排序字段")
	}
	fmt.Println(FmtTracks(os.Stdout, params[0]))
}

// Task9 to Task8
func (e *Exercises) Task9() {
	e.Task8()
}

// Task10 sort.Interface类型也可以适用在其它地方。
// 编写一个IsPalindrome(s sort.Interface) bool函数表明序列s是否是回文序列，
// 换句话说反向排序不会改变这个序列。假设如果!s.Less(i, j) && !s.Less(j, i)则索引i和j上的元素相等。
// go run main.go -ch=7 -task=10
func (e *Exercises) Task10() {
	elems := []int{1, 2, 3, 4, 5, 6, 5, 4, 3, 2, 1}
	fmt.Println(elems, " is palindrome？ ", IsPalindrome(sort.IntSlice(elems)))
	elems = []int{1, 2, 3, 4, 5, 6, 6, 5, 4, 3, 2, 1}
	fmt.Println(elems, " is palindrome？ ", IsPalindrome(sort.IntSlice(elems)))
	elems = []int{1, 2, 3, 4, 5, 6, 4, 3, 2, 1}
	fmt.Println(elems, " is palindrome？ ", IsPalindrome(sort.IntSlice(elems)))
}

// Task11
// 1. 增加额外的handler让客户端可以创建，读取，更新和删除数据库记录。例如，一个形如 /update?item=socks&price=6 的请求会更新库存清单里一个货品的价格并且当这个货品不存在或价格无效时返回一个错误值。（注意：这个修改会引入变量同时更新的问题）
// 2. 修改/list的handler让它把输出打印成一个HTML的表格而不是文本。html/template包（§4.6）可能会对你有帮助。
// go run main.go -ch=7 -task=11
func (e *Exercises) Task11() {
	DB.list(os.Stdout)
}

// Task12 to Task11
func (e *Exercises) Task12() {
	e.Task11()
}

// Task13
// 1. 为Expr增加一个String方法来打印美观的语法树。当再一次解析的时候，检查它的结果是否生成相同的语法树。
// 2. 定义一个新的满足Expr接口的具体类型并且提供一个新的操作例如对它运算单元中的最小值的计算。因为Parse函数不会创建这个新类型的实例，为了使用它你可能需要直接构造一个语法树（或者继承parser接口）。
// go run main.go -ch=7 -task=13
func (e *Exercises) Task13() {
	tests := []struct {
		expr string
		env  eval.Env
		want string
	}{
		{"sqrt(A / pi)", eval.Env{"A": 87616, "pi": math.Pi}, "167"},
		{"pow(x, 3) + pow(y, 3)", eval.Env{"x": 12, "y": 1}, "1729"},
		{"pow(x, 3) + pow(y, 3)", eval.Env{"x": 9, "y": 10}, "1729"},
		{"5 / 9 * (F - 32)", eval.Env{"F": -40}, "-40"},
		{"max(sqrt(A / pi), pow(x, 3) + pow(y, 3), 5 / 9 * (F - 32))", eval.Env{"A": 87616, "pi": math.Pi, "x": 12, "y": 1, "F": -40}, "1729"},
	}
	fun := func(env eval.Env, input string) (float64, error) {
		parse, err := eval.Parse(input)
		if err != nil {
			return 0, err
		}
		parse1, err := eval.Parse(parse.String(env))
		if err != nil {
			return 0, err
		}
		fmt.Println(parse.String(env))
		if parse.Eval(env) != parse1.Eval(env) {
			return 0, errors.New(parse.String(env) + "和" + input + "结果不同")
		}
		return parse.Eval(env), nil
	}
	for _, test := range tests {
		if ret, err := fun(test.env, test.expr); err != nil {
			fmt.Println(err)
			return
		} else if test.want != strconv.Itoa(int(ret)) {
			fmt.Println(test.want, test.expr)
			return
		} else {
			fmt.Println("success", ret)
		}
	}
}

// Task14 to Task13
func (e *Exercises) Task14() {
	e.Task13()
}

// Task15 编写一个从标准输入中读取一个单一表达式的程序，用户及时地提供对于任意变量的值，然后在结果环境变量中计算表达式的值。优雅的处理所有遇到的错误。
// go run main.go -ch=7 -task=15
func (e *Exercises) Task15() {
	for {
		var input string
		var err error
		var env map[eval.Var]float64
		var expr eval.Expr

		buf := bufio.NewReader(os.Stdin)
		if err != nil {
			return
		}
		fmt.Printf("请输入表达式：")
		if input, err = buf.ReadString('\n'); err != nil {
			fmt.Println(err)
			return
		}
		if expr, env, err = eval.ParseAndCheck(input); err != nil {
			return
		}
		for key := range env {
			var val float64
			fmt.Printf("请输%s的值：", key)
			if _, err = fmt.Scan(&val); err != nil {
				return
			}
			env[key] = val
		}
		fmt.Println(expr.Eval(env))
	}
}

// Task16 编写一个基于web的计算器程序。
// go run main.go -ch=7 -task=16
// or go run main.go -http=8080 后请求 http://localhost:8080/ch7/eval
func (e *Exercises) Task16() {
	fmt.Println(eval.Run(os.Stdout, "3 + 4"))
}

// Task17  扩展xmlselect程序以便让元素不仅可以通过名称选择，也可以通过它们CSS风格的属性进行选择。可以通过匹配id或者class，同时还有它的名称来进行选择。
// go run main.go -ch=7 -task=17
func (e *Exercises) Task17() {
	file, err := util.GetFile(1, filepath.Join("file", "text.html"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	fmt.Println(XMLSelect(file, "class", "title"))
}

// Task18 使用基于标记的解码API，编写一个可以读取任意XML文档并构造这个文档所代表的通用节点树的程序。节点有两种类型：CharData节点表示文本字符串，和 Element节点表示被命名的元素和它们的属性。每一个元素节点有一个子节点的切片。
// go run main.go -ch=7 -task=18
func (e *Exercises) Task18() {
	file, err := util.GetFile(1, filepath.Join("file", "text.html"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	ret, err := ParseTree(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	indent, err := json.MarshalIndent(ret, "", " ")
	if err != nil {
		return
	}
	fmt.Println(string(indent))
}
