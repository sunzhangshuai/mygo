package ch1

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// globalCycles 全局周期
var globalCycles float64 = 5

// Exercises 练习
type Exercises struct {
}

// palette 调色板
var palette = []color.Color{color.White, color.Black}

// echo 程序
func echo() {
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}

// dup 程序
func dup() {
	counts := make(map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		dupCountLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			dupCountLines(f, counts)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

// lissajous 利萨如图形
func lissajous(out io.Writer, getColorIndex func() uint8) {
	// 周期
	cycles := globalCycles
	const (
		// number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)

	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				getColorIndex())
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

// fetch 获取url
func fetch() {
	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
		fmt.Printf("%s", b)
	}
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
	resp.Body.Close() // don't leak resources
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
func (e *Exercises) Task1() {
	var s, sep string
	for i := 0; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}

// Task2 修改echo程序，使其打印每个参数的索引和值，每个一行。
func (e *Exercises) Task2() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Println(i, os.Args[i])
	}
}

// Task3 做实验测量潜在低效的版本和使用了strings.Join的版本的运行时间差异。
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
func (e *Exercises) Task4() {
	counts := make(map[string]int)
	files := os.Args[5:]
	if len(files) == 0 {
		task4CountLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			task4CountLines(f, counts)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

// Task5 修改前面的Lissajous程序里的调色板，由黑色改为绿色。我们可以用color.RGBA{0xRR, 0xGG, 0xBB, 0xff}来得到#RRGGBB这个色值，三个十六进制的字符串分别代表红、绿、蓝像素。
func (e *Exercises) Task5() {
	// 设置随机数种子
	rand.Seed(time.Now().UTC().UnixNano())
	// 创建文件
	fileFd, _ := os.Create("ch1/output/task5_" + strconv.Itoa(time.Now().Hour()) + strconv.Itoa(time.Now().Minute()) + strconv.Itoa(time.Now().Second()) + ".gif")
	// 修改调色板
	palette = []color.Color{color.White, color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xff}}
	getColorIndex := func() uint8 {
		return 1
	}
	lissajous(fileFd, getColorIndex)
}

// Task6 修改Lissajous程序，修改其调色板来生成更丰富的颜色，然后修改SetColorIndex的第三个参数，看看显示结果吧。
func (e *Exercises) Task6() {
	// 设置随机数种子
	rand.Seed(time.Now().UTC().UnixNano())
	// 创建文件
	fileFd, _ := os.Create("ch1/output/task6_" + strconv.Itoa(time.Now().Hour()) + strconv.Itoa(time.Now().Minute()) + strconv.Itoa(time.Now().Second()) + ".gif")
	// 修改调色板
	palette = []color.Color{color.Black}
	for i := 0; i < 100; i++ {
		palette = append(palette, color.RGBA{R: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)), B: uint8(rand.Intn(255)), A: 0xff})
	}
	var colorIndex uint8 = 0
	getColorIndex := func() uint8 {
		colorIndex = uint8((int(colorIndex)+1)%len(palette)-1) + 1
		return colorIndex
	}
	lissajous(fileFd, getColorIndex)
}

// Task7 函数调用io.Copy(dst, src)会从src中读取内容，并将读到的结果写入到dst中，使用这个函数替代掉例子中的ioutil.ReadAll来拷贝响应结构体到os.Stdout，避免申请一个缓冲区（例子中的b）来存储。记得处理io.Copy返回结果中的错误。
func (e *Exercises) Task7() {
	for _, url := range os.Args[5:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		_, err = io.Copy(os.Stdout, resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
	}
}

// Task8 修改fetch这个范例，如果输入的url参数没有 http:// 前缀的话，为这个url加上该前缀。你可能会用到strings.HasPrefix这个函数。
func (e *Exercises) Task8() {
	for _, url := range os.Args[5:] {
		if !strings.HasPrefix(url, "http://") {
			url = "http://" + url
		}
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		_, err = io.Copy(os.Stdout, resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
	}
}

// Task9 修改fetch打印出HTTP协议的状态码，可以从resp.Status变量得到该状态码。
func (e *Exercises) Task9() {
	for _, url := range os.Args[5:] {
		if !strings.HasPrefix(url, "http://") {
			url = "http://" + url
		}
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "status: %v\n", resp.Status)
		resp.Body.Close()
	}
}

// Task10 找一个数据量比较大的网站，用本小节中的程序调研网站的缓存策略，对每个URL执行两遍请求，
// 查看两次时间是否有较大的差别，并且每次获取到的响应内容是否一致，修改本节中的程序，将响应结果输出，以便于进行对比。
func (e *Exercises) Task10() {
	fetchAll([]string{"http://www.baidu.com", "http://www.baidu.com"})
}

// Task11 在fetchall中尝试使用长一些的参数列表，比如使用在alexa.com的上百万网站里排名靠前的。如果一个网站没有回应，程序将采取怎样的行为？（Section8.9 描述了在这种情况下的应对机制）。
func (e *Exercises) Task11() {
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

// Task12 修改Lissajour服务，从URL读取变量，比如你可以访问 http://localhost:8000/?cycles=20 这个URL，这样访问可以将程序里的cycles默认的5修改为20。字符串转换为数字可以调用strconv.Atoi函数。你可以在godoc里查看strconv.Atoi的详细说明。
func (e *Exercises) Task12() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		var err error
		if err = request.ParseForm(); err != nil {
			log.Print(err)
		}
		// 获取参数
		cycles := request.Form.Get("cycles")
		log.Print(cycles)
		if globalCycles, err = strconv.ParseFloat(cycles, 64); err != nil {
			globalCycles = 5
			log.Print(err)
		}
		lissajous(writer, func() uint8 {
			return 1
		})
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
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
	fmt.Printf("长度为%d：\n", len(stringList))
	fmt.Printf("不使用join的时间%d：\n", t2.Sub(t1))
	fmt.Printf("使用join的时间%d：\n", t3.Sub(t2))
}

// task4CountLines 统计字符行数
func task4CountLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		if count, ok := counts[input.Text()]; ok {
			fmt.Printf("文本【%s】第【%d】出现", input.Text(), count+1)
		}
		counts[input.Text()]++
	}
	// NOTE: ignoring potential errors from input.Err()
}
