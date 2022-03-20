package ch4

import (
	"bufio"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gopl.io/ch4/github"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"exercises/ch2"
	"exercises/ch4/gitissue"
	"exercises/util"
)

type XKCD struct {
	Img   string
	Title string
}

// Exercises 练习
type Exercises struct {
}

// Task1 编写一个函数，计算两个SHA256哈希码中不同bit的数目。
// go run main.go -ch=4 -task=1 张三 李四
func (e *Exercises) Task1() {
	params := util.GetParam()
	if len(params) < 2 {
		fmt.Println("参数不足")
		return
	}
	c1 := sha256.Sum256([]byte(params[0]))
	c2 := sha256.Sum256([]byte(params[1]))
	fmt.Println(hex.EncodeToString(c1[:]))
	fmt.Println(hex.EncodeToString(c2[:]))

	ret := 0
	for i, val := range c1 {
		ret += ch2.PopCount(uint64(val | c2[i]))
	}
	fmt.Println(ret)
}

// Task2 编写一个程序，默认情况下打印标准输入的SHA256编码，并支持通过命令行flag定制，输出SHA384或SHA512哈希算法。
// go run main.go -ch=4 -task=2 张三 sha256
func (e *Exercises) Task2() {
	params := util.GetParam()
	h := "sha256"
	if len(params) < 1 {
		fmt.Println("参数不足")
		return
	} else if len(params) > 1 {
		h = params[1]
	}

	switch h {
	case "sha256":
		sha := sha256.Sum256([]byte(params[0]))
		fmt.Println(hex.EncodeToString(sha[:]))
	case "sha384":
		sha := sha512.Sum384([]byte(params[0]))
		fmt.Println(hex.EncodeToString(sha[:]))
	case "sha512":
		sha := sha512.Sum512([]byte(params[0]))
		fmt.Println(hex.EncodeToString(sha[:]))
	}
}

// Task3  重写reverse函数，使用数组指针代替slice。
// go run main.go -ch=4 -task=3
func (e *Exercises) Task3() {
	s := [...]int{1, 2, 3, 4, 5, 6, 7}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Task4 编写一个rotate函数，通过一次循环完成旋转。
// go run main.go -ch=4 -task=4
func (e *Exercises) Task4() {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	rotate(s, 1)
}

// Task5  写一个函数在原地完成消除[]string中相邻重复的字符串的操作。
// go run main.go -ch=4 -task=5
func (e *Exercises) Task5() {
	arr := []string{"a", "b", "b", "b", "c", "c", "d"}
	var slow, fast int

	for fast = 1; fast < len(arr); fast++ {
		if arr[fast] != arr[slow] {
			slow++
			arr[slow] = arr[fast]
		}
	}

	fmt.Println(arr[:slow+1])
}

// Task6 编写一个函数，原地将一个UTF-8编码的[]byte类型的slice中相邻的空格（参考unicode.IsSpace）替换成一个空格返回。
// go run main.go -ch=4 -task=6
func (e *Exercises) Task6() {
	str := "this is go   张帅,  but i  love"
	byt := []byte(str)

	var slow, fast int
	for _, fast = utf8.DecodeRuneInString(string(byt)); fast < len(string(byt)); {
		rF, falseSize := utf8.DecodeRuneInString(string(byt[fast:]))
		rS, slowSize := utf8.DecodeRuneInString(string(byt[slow:]))
		if !unicode.IsSpace(rF) || rF != rS {
			slow += slowSize
			copy(byt[slow:slow+falseSize], byt[fast:fast+falseSize])
		}
		fast += falseSize
	}
	fmt.Println(string(byt[:slow+1]))
}

// Task7 修改reverse函数用于原地反转UTF-8编码的[]byte。是否可以不用分配额外的内存？
func (e *Exercises) Task7() {
	str := "丑ـبــثわたし 张帅"
	byt := []byte(str)

	for l := len(byt); l > 0; {
		r, size := utf8.DecodeRuneInString(string(byt[0:]))
		copy(byt[0:l], byt[0+size:l])
		copy(byt[l-size:l], string(r))
		l -= size
	}
	fmt.Println(string(byt))
}

// Task8 修改charcount程序，使用unicode.IsLetter等相关的函数，统计字母、数字等Unicode中不同的字符类别。
// go run main.go -ch=4 -task=8
func (e *Exercises) Task8() {
	counts := make(map[rune]int)    // counts of Unicode characters
	var utflen [utf8.UTFMax + 1]int // count of lengths of UTF-8 encodings
	invalid := 0                    // count of invalid UTF-8 characters
	letterCount, numberCount, otherCount, characters := 0, 0, 0, 0
	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}

		if unicode.IsNumber(r) {
			numberCount++
			d, _ := strconv.Atoi(string(r))
			fmt.Printf("%d is number\n", d)
		} else if r < 128 {
			letterCount++
			fmt.Printf("%v is letter\n", string(r))
		} else if unicode.IsLetter(r) {
			characters++
			fmt.Printf("%v is characters\n", string(r))
		} else {
			otherCount++
		}

		counts[r]++
		utflen[n]++
	}
	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}
	fmt.Print("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}

	fmt.Printf("characters： %v,\t letter: %v,\t number,%v\t other: %v", characters, letterCount, numberCount, otherCount)
}

// Task9  编写一个程序wordfreq程序，报告输入文本中每个单词出现的频率。
// 在第一次调用Scan前先调用input.Split(bufio.ScanWords)函数，这样可以按单词而不是按行输入。
// go run main.go -ch=4 -task=9
func (e *Exercises) Task9() {
	counts := make(map[string]int)

	input := bufio.NewScanner(os.Stdin)
	input.Split(bufio.ScanWords)

	for input.Scan() {
		word := input.Text()
		if _, ok := counts[word]; !ok {
			counts[word] = 0
		}
		counts[word]++
	}
	fmt.Printf("%#v", counts)
}

// Task10 修改issues程序，根据问题的时间进行分类，比如不到一个月的、不到一年的、超过一年。
// go run main.go -ch=4 -task=10
func (e *Exercises) Task10() {
	var err error
	var gitResult *github.IssuesSearchResult
	var jsonRes []byte

	if gitResult, err = github.SearchIssues([]string{"repo:golang/go", "is:open", "json decoder"}); err != nil {
		log.Fatal(err)
	}

	mapRes := make(map[string][]string)
	mapRes["far"] = make([]string, 0)
	mapRes["middle"] = make([]string, 0)
	mapRes["near"] = make([]string, 0)
	nowTime := time.Now()
	for _, item := range gitResult.Items {
		createdAt := item.CreatedAt
		if createdAt.AddDate(1, 0, 0).Before(nowTime) {
			mapRes["far"] = append(mapRes["far"], item.HTMLURL)
		} else if createdAt.AddDate(0, 1, 0).Before(nowTime) {
			mapRes["middle"] = append(mapRes["middle"], item.HTMLURL)
		} else {
			mapRes["near"] = append(mapRes["near"], item.HTMLURL)
		}
	}

	if jsonRes, err = json.MarshalIndent(mapRes, "", "\t"); err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonRes))
}

// Task11 编写一个工具，允许用户在命令行创建、读取、更新和关闭GitHub上的issue，当必要的时候自动打开用户默认的编辑器用于输入文本信息。
// go run main.go -ch=4 -task=11
func (e *Exercises) Task11() {
	fmt.Println(gitissue.Run())
}

// Task12 流行的web漫画服务xkcd也提供了JSON接口。
// 例如，一个 https://xkcd.com/571/info.0.json 请求将返回一个很多人喜爱的571编号的详细描述。
// 下载每个链接（只下载一次）然后创建一个离线索引。编写一个xkcd工具，使用这些离线索引，打印和命令行输入的检索词相匹配的漫画的URL。
// go run main.go -ch=4 -task=12
func (e *Exercises) Task12() {
	var err error
	var file *os.File

	if file, err = util.NewFile(1, filepath.Join("file", "xkcd.txt")); err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	result := make(map[string]string)
	for i := 1; i < 20; i++ {
		uri := "https://xkcd.com/" + strconv.Itoa(i) + "/info.0.json"

		// 处理每条数据
		if func() error {
			res, err := http.Get(uri)
			if err != nil {
				return err
			}
			defer res.Body.Close()

			var item XKCD
			item2 := make(map[string]interface{})

			json.NewDecoder(res.Body).Decode(&item2)
			item.Title = item2["title"].(string)
			item.Img = item2["img"].(string)
			result[item.Title] = item.Img
			var byt []byte
			if byt, err = json.MarshalIndent(item2, "", "\t"); err != nil {
				return err
			}
			str := string(byt) + "\n"
			if _, err = file.WriteString(str); err != nil {
				return err
			}
			return nil
		}() != nil {
			continue
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("请输入要查找的图片")
	for scanner.Scan() {
		str := scanner.Text()
		for title, img := range result {
			if strings.Index(title, str) != -1 {
				fmt.Println("title:"+title, "img:"+img)
			}
		}
	}
}

// Task13 使用开放电影数据库的JSON服务接口，允许你检索和下载 https://omdbapi.com/ 上电影的名字和对应的海报图像。
// 编写一个poster工具，通过命令行输入的电影名字，下载对应的海报。
func (e *Exercises) Task13() {
	// 需要注册，不做了
}

// Task14 创建一个web服务器，查询一次GitHub，然后生成BUG报告、里程碑和对应的用户信息。
// go run main.go -ch=4 -task=14
// go run main.go -http=8080 后请求 http://localhost:8080/ch4/template
func (e *Exercises) Task14() {
	var err error
	var file *os.File

	if file, err = util.NewFile(1, filepath.Join("file", "template.html")); err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	RunTemplate(file)
}

// rotate 旋转
func rotate(s []int, l int) {
	fmt.Println(s)
	zero := s[0]
	sLen := len(s)

	next := func(i int) int {
		return (i + sLen - l) % sLen
	}

	var i int
	for i = 0; next(i) != 0; i = next(i) {
		fmt.Println(i, next(i))
		s[i] = s[next(i)]
	}
	s[i] = zero
	fmt.Println(s)
}
