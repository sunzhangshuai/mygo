// 编写一个函数，原地将一个UTF-8编码的[]byte类型的slice中相邻的空格（参考unicode.IsSpace）替换成一个空格返回
package main

import (
	"fmt"
	"unicode"
)

func main() {
	str := "this is go   test,  but i  love"
	bStr := []byte(str)
	for i := 1; i < len(bStr); i++ {
		j := i
		for ; j < len(bStr); j++ {
			if !unicode.IsSpace(rune(bStr[j])) || !unicode.IsSpace(rune(bStr[j-1])) {
				break
			}
		}
		if i != j {
			copy(bStr[i:], bStr[j:])
			bStr = bStr[:len(bStr)-j+i]
		}
	}
	fmt.Println(string(bStr))
}
