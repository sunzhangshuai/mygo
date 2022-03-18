package ch4

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
	"unicode"
	"unicode/utf8"

	"exercises/ch2"
)

// Exercises 练习
type Exercises struct {
}

// Task1 编写一个函数，计算两个SHA256哈希码中不同bit的数目。
// go run main.go -ch 4 -task 1 张三 李四
func (e *Exercises) Task1() {
	c1 := sha256.Sum256([]byte(os.Args[5]))
	c2 := sha256.Sum256([]byte(os.Args[6]))
	fmt.Println(hex.EncodeToString(c1[:]))
	fmt.Println(hex.EncodeToString(c2[:]))

	ret := 0
	for i, val := range c1 {
		ret += ch2.PopCount(uint64(val | c2[i]))
	}
	fmt.Println(ret)
}

// Task2 编写一个程序，默认情况下打印标准输入的SHA256编码，并支持通过命令行flag定制，输出SHA384或SHA512哈希算法。
// go run main.go -ch 4 -task 2 张三 sha256
func (e *Exercises) Task2() {
	h := "sha256"
	if len(os.Args) > 6 {
		h = os.Args[6]
	}

	switch h {
	case "sha256":
		sha := sha256.Sum256([]byte(os.Args[5]))
		fmt.Println(hex.EncodeToString(sha[:]))
	case "sha384":
		sha := sha512.Sum384([]byte(os.Args[5]))
		fmt.Println(hex.EncodeToString(sha[:]))
	case "sha512":
		sha := sha512.Sum512([]byte(os.Args[5]))
		fmt.Println(hex.EncodeToString(sha[:]))
	}
}

// Task3  重写reverse函数，使用数组指针代替slice。
// go run main.go -ch 4 -task 3
func (e *Exercises) Task3() {
	s := [...]int{1, 2, 3, 4, 5, 6, 7}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Task4 编写一个rotate函数，通过一次循环完成旋转。
// go run main.go -ch 4 -task 4
func (e *Exercises) Task4() {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	rotate(s, 1)
}

// Task5  写一个函数在原地完成消除[]string中相邻重复的字符串的操作。
// go run main.go -ch 4 -task 5
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
// go run main.go -ch 4 -task 6
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

func (e *Exercises) Task8() {

}

func (e *Exercises) Task9() {

}

func (e *Exercises) Task10() {

}

func (e *Exercises) Task11() {

}

func (e *Exercises) Task12() {

}

func (e *Exercises) Task13() {

}

func (e *Exercises) Task14() {

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
