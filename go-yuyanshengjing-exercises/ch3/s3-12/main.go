// 编写一个函数，判断两个字符串是否是相互打乱的，也就是说它们有着相同的字符，但是对应不同的顺序。
package main

import (
	"fmt"
)

func main() {
	fmt.Println(compare("aaa", "aab"))
	fmt.Println(compare("aaa张三", "aab张三"))
	fmt.Println(compare("aaa张三", "张三aaa"))
}

func compare(a string, b string) bool {
	mapA := make(map[rune]int)
	for _, s := range a{
		mapA[s]++
	}
	for _, s := range b{
		mapA[s]--
	}
	for _, val :=  range mapA{
		if val != 0 {
			return false;
		}
	}
	return true
}
