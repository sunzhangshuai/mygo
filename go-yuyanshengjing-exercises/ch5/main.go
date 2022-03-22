package ch5

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"exercises/ch3"
	"exercises/util"
)

// Exercises 练习
type Exercises struct {
}

// Task1
// 1. 修改findlinks代码中遍历n.FirstChild链表的部分，将循环调用visit，改成递归调用。
// 2. 编写函数，记录在HTML树中出现的同名元素的次数。
// 3. 编写函数输出所有text结点的内容。注意不要访问<script>和<style>元素，因为这些元素对浏览者是不可见的。
// 4. 扩展visit函数，使其能够处理其他类型的结点，如images、scripts和style sheets。
// go run main.go -ch=5 -task=1
func (e *Exercises) Task1() {
	var doc *html.Node
	var err error
	var ret *http.Response

	if ret, err = http.Get("http://gopl.io/"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer ret.Body.Close()

	if doc, err = html.Parse(ret.Body); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	visit(doc)
	printData()
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

// Task5 实现countWordsAndImages。
// go run main.go -ch=5 -task=5 http://gopl.io
func (e *Exercises) Task5() {
	params := util.GetParam()
	if len(params) == 0 {
		fmt.Println("参数不足")
		os.Exit(1)
	}
	fmt.Println(CountWordsAndImages(params[0]))
	return
}

// Task6 修改gopl.io/ch3/surface（§3.2）中的corner函数，将返回值命名，并使用bare return。
// go run main.go -ch=5 -task=6 sunchen
func (e *Exercises) Task6() {
	var err error
	var file *os.File
	params := util.GetParam()
	if len(params) == 0 {
		fmt.Println("需要传入文件名")
		return
	}
	if file, err = util.NewFile(1, filepath.Join("file", params[0]+".svg")); err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	ch3.Surface(file, "")
}

// Task7
// 1. 完善startElement和endElement函数，使其成为通用的HTML输出器。 要求：输出注释结点，文本结点以及每个元素的属性（< a href='...'>）。使用简略格式输出没有孩子结点的元素（即用<img/>代替<img></img>）。
// 2. 修改pre和post函数，使其返回布尔类型的返回值。返回false时，中止forEachNoded的遍历。使用修改后的代码编写ElementByID函数，根据用户输入的id查找第一个拥有该id元素的HTML元素，查找成功后，停止遍历。
// 3. gopl.io/ch5/outline2（5.5节）的startElement和endElement共用了全局变量depth，将它们修改为匿名函数，使其共享outline中的局部变量。
// go run main.go -ch=5 -task=7 http://gopl.io
func (e *Exercises) Task7() {
	var doc *html.Node
	var err error
	var ret *http.Response

	params := util.GetParam()
	if len(params) == 0 {
		fmt.Println("参数不足")
		os.Exit(1)
	}

	if ret, err = http.Get(params[0]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer ret.Body.Close()

	if doc, err = html.Parse(ret.Body); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	var depth int

	depthF := func(action string) int {
		switch action {
		case "add":
			res := depth
			depth++
			return res
		case "sub":
			depth--
			return depth
		}
		return depth
	}

	forEachNode(doc, startElement, endElement, func(ID string) bool {
		return ID == "2"
	}, depthF)
}

// Task8 to Task7
func (e *Exercises) Task8() {
	e.Task7()
}

// Task9 编写函数expand，将s中的"foo"替换为f("foo")的返回值。
// go run main.go -ch=5 -task=9 sadasdafooasdasfoofoofoo
func (e *Exercises) Task9() {
	params := util.GetParam()
	if len(params) == 0 {
		fmt.Println("参数不足")
		os.Exit(1)
	}
	fmt.Println(expand(params[0], func(s string) string {
		return fmt.Sprintf("[[[%s]]]", s)
	}))
}

// Task10
// 1. 重写topoSort函数，用map代替切片并移除对key的排序代码。验证结果的正确性（结果不唯一）。
// 2. 现在线性代数的老师把微积分设为了前置课程。完善topSort，使其能检测有向图中的环。 todo
// go run main.go -ch=5 -task=10
func (e *Exercises) Task10() {
	res, err := json.MarshalIndent(topologySort(preReqs), "", "\t")
	if err != nil {
		fmt.Println(res)
		os.Exit(1)
	}
	fmt.Println(string(res))
}

// Task11 to Task10
func (e *Exercises) Task11() {
	e.Task10()
}

// Task12 to Task7
func (e *Exercises) Task12() {
	e.Task7()
}

// Task13 修改crawl，使其能保存发现的页面，必要时，可以创建目录来保存这些页面。只保存来自原始域名下的页面。假设初始页面在golang.org下，就不要保存vimeo.com下的页面。
func (e *Exercises) Task13() {}

// Task14 使用breadthFirst遍历其他数据结构。比如，topoSort例子中的课程依赖关系（有向图）、个人计算机的文件层次结构（树）；你所在城市的公交或地铁线路（无向图）。
func (e *Exercises) Task14() {}

// expand 扩展
func expand(s string, f func(string) string) string {
	return strings.Replace(s, "foo", f("foo"), -1)
}
