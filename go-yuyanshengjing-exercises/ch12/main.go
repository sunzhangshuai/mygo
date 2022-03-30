package ch12

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"exercises/ch12/decode"
	"exercises/ch12/encode"
)

// Exercises 练习
type Exercises struct {
}

var strangelove Movie

var test, test1 map[interface{}]interface{}

// Task1
// 1. 扩展Display函数，使它可以显示包含以结构体或数组作为map的key类型的值。
// 2. 增强display函数的稳健性，通过记录边界的步数来确保在超出一定限制后放弃递归。（在13.3节，我们会看到另一种探测数据结构是否存在环的技术。）
// go run main -ch=12 -task=1
func (e *Exercises) Task1() {
	Display("strangelove", strangelove, 2)
	Display("test", test, 9)
}

// Task2 to Task1
func (e *Exercises) Task2() {
	e.Task1()
}

// Task3 实现encode函数缺少的分支。将布尔类型编码为t和nil，浮点数编码为Go语言的格式，复数1+2i编码为#C(1.0 2.0)格式。
// 接口编码为类型名和值对，例如（"[]int" (1 2 3)），但是这个形式可能会造成歧义：reflect.Type.String方法对于不同的类型可能返回相同的结果。
// go run main.go -ch=12 -task=3
func (e *Exercises) Task3() {
	str, _ := encode.Marshal(strangelove)
	fmt.Println(string(str))
	str, _ = encode.Marshal(test1)
	fmt.Println(string(str))
}

// Task4 修改encode函数，以上面的格式化形式输出S表达式。
// go run main.go -ch=12 -task=4
func (e *Exercises) Task4() {
	str, _ := encode.MarshalIndent(strangelove)
	fmt.Println(string(str))
	str, _ = encode.MarshalIndent(test1)
	fmt.Println(string(str))
}

// Task5 修改encode函数，用JSON格式代替S表达式格式。然后使用标准库提供的json.Unmarshal解码器来验证函数是正确的。
// go run main.go -ch=12 -task=5
func (e *Exercises) Task5() {
	strangelove.Actor = nil
	str, _ := encode.MarshalJson(strangelove)
	sStr := string(str)
	fmt.Println(sStr)
	str1, _ := json.Marshal(strangelove)
	sStr1 := string(str1)
	fmt.Println(sStr1)

	fmt.Println(strings.Compare(string(str1), string(str)))
	str, _ = encode.MarshalJson(test1)
	fmt.Println(string(str))
}

// Task6 修改encode，作为一个优化，忽略对是零值对象的编码。
func (e *Exercises) Task6() {
	e.Task3()
}

// Task7 创建一个基于流式的API，用于S表达式的解码，和json.Decoder(§4.5)函数功能类似。
// go run main.go -ch=12 -task=7
func (e *Exercises) Task7() {
	str, _ := encode.Marshal(strangelove)
	movie := Movie{}
	if err := decode.Unmarshal(str, &movie); err != nil {
		fmt.Println(err)
		return
	}
	str1, _ := json.Marshal(movie)
	sStr1 := string(str1)
	fmt.Println(sStr1)
}

// Task8
// 1. sexpr.Unmarshal函数和json.Unmarshal一样，都要求在解码前输入完整的字节slice。 定义一个和json.Decoder类似的sexpr.Decoder类型，支持从一个io.Reader流解码。修改sexpr.Unmarshal函数，使用这个新的类型实现。
// 2. 修改S表达式的编码器（§12.4）和解码器（§12.6），采用和encoding/json包（§4.5）类似的方式使用成员标签中的sexpr:"..."字串。
// go run main.go -ch=12 -task=8
func (e *Exercises) Task8() {
	str, _ := encode.Marshal(strangelove)
	movie := Movie{}
	if err := decode.NewDecoder(bytes.NewReader(str)).Decode(&movie); err != nil {
		fmt.Println(err)
		return
	}
	str1, _ := json.Marshal(movie)
	sStr1 := string(str1)
	fmt.Println(sStr1)
}

// Task9  编写一个基于标记的API用于解码S表达式，参考xml.Decoder（7.14）的风格。
// 你将需要五种类型的标记：Symbol、String、Int、StartList和EndList。
func (e *Exercises) Task9() {
	// 不做了
}

// Task10 扩展sexpr.Unmarshal函数，支持布尔型、浮点数和interface类型的解码，使用 练习 12.3： 的方案。
//（提示：要解码接口，你需要将name映射到每个支持类型的reflect.Type。）
// go run main.go -ch=12 -task=10
func (e *Exercises) Task10() {
	strangelove.Inter = 1234
	str, _ := encode.Marshal(strangelove)
	movie := Movie{}
	if err := decode.NewDecoder(bytes.NewReader(str)).Decode(&movie); err != nil {
		fmt.Println(err)
		return
	}
	str1, _ := json.Marshal(movie)
	sStr1 := string(str1)
	fmt.Println(sStr1)
}

// Task11 编写相应的Pack函数，给定一个结构体值，Pack函数将返回合并了所有结构体成员和值的URL。
// go run main.go -ch=12 -task=11
func (e *Exercises) Task11() {
	type data struct {
		Labels     []string `http:"l"`
		MaxResults int      `http:"max"`
		Exact      bool     `http:"x"`
	}
	d := data{
		Labels:     []string{"golang", "programming"},
		MaxResults: 20,
		Exact:      true,
	}
	fmt.Println(Pack(&d))
}

// Task12 扩展成员标签以表示一个请求参数的有效值规则。例如，一个字符串可以是有效的email地址或一个信用卡号码，还有一个整数可能需要是有效的邮政编码。修改Unpack函数以检查这些规则。
// go run main.go -http=8080 后请求 http://localhost:8080/ch12/unpack?l=golang&l=programming&max=100
// go run main.go -ch=12 -task=12
func (e *Exercises) Task12() {
	ret, err := http.Get("http://localhost:8080/ch12/unpack?l=golang&l=programming&max=100")
	if err != nil {
		return
	}
	defer ret.Body.Close()
	io.Copy(os.Stdout, ret.Body)
}

// Task13 to Task8
func (e *Exercises) Task13() {
	e.Task8()
}

func init() {
	strangelove = Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Color:    true,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},

		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}

	key1 := [4]int{1, 2, 3, 4}
	type data struct {
		Key   string
		Value string
		Next  *data
	}
	key2 := data{"key2", "key2", nil}
	key3 := data{"key3", "key3", &key2}
	key4 := data{"key2", "key2", nil}
	key2.Next = &key3
	test = map[interface{}]interface{}{
		key1: key2,
	}

	test1 = map[interface{}]interface{}{
		key1: key4,
		key4: key1,
	}
}
