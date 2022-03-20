package ch1

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"exercises/util"
)

// Exercises 练习
type Exercises struct {
}

// fetchAll 并发获取url
func fetchAll(urls []string) {
	start := time.Now()
	ch := make(chan string)
	for _, url := range urls {
		if !strings.HasPrefix(url, "http://") {
			url = "http://" + url
		}
		go fetchOne(url, ch) // start a goroutine
	}
	for range urls {
		fmt.Println(<-ch) // receive from channel ch
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

// fetchOne 获取1个url
func fetchOne(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}
	nBytes, err := io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nBytes, url)
}

// dupCountLines 统计字符行数
func dupCountLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[input.Text()]++
	}
	// NOTE: ignoring potential errors from input.Err()
}

// Task1 修改echo程序，使其能够打印os.Args[0]，即被执行命令本身的名字。
// 1. 修改echo程序，使其能够打印os.Args[0]，即被执行命令本身的名字。
// 2. 修改echo程序，使其打印每个参数的索引和值，每个一行。
// go run main.go -ch=1 -task=1
func (e *Exercises) Task1() {
	var s, sep string
	for i := 0; i < len(os.Args); i++ {
		fmt.Println(i, os.Args[i])
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}

// Task2 to Task1
func (e *Exercises) Task2() {
	e.Task1()
}

// Task3 做实验测量潜在低效的版本和使用了strings.Join的版本的运行时间差异。
// go run main.go -ch=1 -task=3
func (e *Exercises) Task3() {
	stringList := make([]string, 10)
	for i := 0; i < 10; i++ {
		stringList[i] = "str" + strconv.Itoa(i)
	}
	task3Test(stringList)

	stringList = make([]string, 1000)
	for i := 0; i < 10; i++ {
		stringList[i] = "str" + strconv.Itoa(i)
	}
	task3Test(stringList)
}

// Task4 修改dup，出现重复的行时打印文件名称。
// go run main.go -ch=1 -task=4
func (e *Exercises) Task4() {
	var err error
	var file *os.File

	counts := make(map[string]int)
	files := util.GetParam()
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	}

	for _, arg := range files {
		func() {
			if file, err = os.Open(arg); err != nil {
				fmt.Printf("dup2: %v\n", err)
				return
			}
			countLines(file, counts)
			defer file.Close()
		}()
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

// Task5
// 1. 修改前面的Lissajous程序里的调色板，由黑色改为绿色。我们可以用color.RGBA{0xRR, 0xGG, 0xBB, 0xff}来得到#RRGGBB这个色值，三个十六进制的字符串分别代表红、绿、蓝像素。
// 2. 修改Lissajous程序，修改其调色板来生成更丰富的颜色，然后修改SetColorIndex的第三个参数，看看显示结果吧。
// 3. 修改Lissajour服务，从URL读取变量，比如你可以访问 http://localhost:8000/?cycles=20 这个URL，这样访问可以将程序里的cycles默认的5修改为20。字符串转换为数字可以调用strconv.Atoi函数。你可以在godoc里查看strconv.Atoi的详细说明。
// go run main.go -ch=1 -task=5 zs
// go run main.go -http=8080 后请求 http://localhost:8080/ch1/lissajous?cycles=20
func (e *Exercises) Task5() {
	var file *os.File
	var err error
	params := util.GetParam()
	if len(params) == 0 {
		fmt.Println("需要输入文件名")
		return
	}

	// 设置随机数种子
	rand.Seed(time.Now().UTC().UnixNano())
	// 创建文件
	if file, err = util.NewFile(1, filepath.Join("file", params[0]+".gif")); err != nil {
		fmt.Println(err)
		return
	}
	// 修改调色板
	for i := 0; i < 100; i++ {
		palette = append(palette, color.RGBA{R: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)), B: uint8(rand.Intn(255)), A: 0xff})
	}
	var colorIndex uint8 = 0
	getColorIndex := func() uint8 {
		colorIndex = uint8((int(colorIndex)+1)%len(palette)-1) + 1
		return colorIndex
	}
	fmt.Println(Lissajous(file, getColorIndex, globalCycles))
}

// Task6 to Task5
func (e *Exercises) Task6() {
	e.Task5()
}

// Task7
// 1. 函数调用io.Copy(dst, src)会从src中读取内容，并将读到的结果写入到dst中， 使用这个函数替代掉例子中的ioutil.ReadAll来拷贝响应结构体到os.Stdout，避免申请一个缓冲区（例子中的b）来存储。 记得处理io.Copy返回结果中的错误。
// 2. 修改fetch这个范例，如果输入的url参数没有 http:// 前缀的话，为这个url加上该前缀。你可能会用到strings.HasPrefix这个函数。
// 3. 修改fetch打印出HTTP协议的状态码，可以从resp.Status变量得到该状态码。
// go run main.go -ch=1 -task=7 baidu.com zhihu.com
func (e *Exercises) Task7() {
	var resp *http.Response
	var err error
	urls := util.GetParam()

	for _, url := range urls {
		func() {
			if !strings.HasPrefix(url, "http://") || !strings.HasPrefix(url, "https://") {
				url = "http://" + url
			}
			if resp, err = http.Get(url); err != nil {
				fmt.Printf("fetch: %v\n", err)
				os.Exit(1)
			}
			defer resp.Body.Close()
			if _, err = io.Copy(os.Stdout, resp.Body); err != nil {
				fmt.Printf("fetch: reading %s: %v\n", url, err)
				os.Exit(1)
			}
			fmt.Printf("status: %v\n", resp.Status)
		}()
	}
}

// Task8 to Task7
func (e *Exercises) Task8() {
	e.Task7()
}

// Task9 to Task7
func (e *Exercises) Task9() {
	e.Task7()
}

// Task10
// 1. 找一个数据量比较大的网站，用本小节中的程序调研网站的缓存策略，对每个URL执行两遍请求， 查看两次时间是否有较大的差别，并且每次获取到的响应内容是否一致，修改本节中的程序，将响应结果输出，以便于进行对比。
// 2. 在fetchall中尝试使用长一些的参数列表，比如使用在alexa.com的上百万网站里排名靠前的。如果一个网站没有回应，程序将采取怎样的行为？（Section8.9 描述了在这种情况下的应对机制）。
// go run main.go -ch=1 -task=10
func (e *Exercises) Task10() {
	urls := []string{
		"www.tmall.com",
		"www.baidu.com",
		"www.qq.com",
		"www.sohu.com",
		"www.taobao.com",
		"www.360.com",
		"www.jd.com",
		"www.weibo.com",
		"www.sina.com",
		"www.zhanqi.com",
	}
	for i := 0; i < 3; i++ {
		urls = append(urls, urls...)
	}
	fetchAll(urls)
}

// Task11 to Task10
func (e *Exercises) Task11() {
	e.Task10()
}

// Task12 to Task5
func (e *Exercises) Task12() {
	e.Task5()
}

// task3Test 任务3的测试
func task3Test(stringList []string) {
	t1 := time.Now()
	for i := 0; i < 1000; i++ {
		var s, sep string
		for i := 0; i < len(stringList); i++ {
			s += sep + stringList[i]
			sep = " "
		}
	}
	t2 := time.Now()
	for i := 0; i < 1000; i++ {
		strings.Join(stringList, " ")
	}
	t3 := time.Now()
	fmt.Printf("长度为：%d\n", len(stringList))
	fmt.Printf("不使用join的时间：%s\n", t2.Sub(t1).String())
	fmt.Printf("使用join的时间：%s\n", t3.Sub(t2).String())
}

// countLines 统计字符行数
func countLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		if count, ok := counts[input.Text()]; ok {
			fmt.Printf("文本【%s】第【%d】出现", input.Text(), count+1)
		}
		counts[input.Text()]++
	}
	// NOTE: ignoring potential errors from input.Err()
}
