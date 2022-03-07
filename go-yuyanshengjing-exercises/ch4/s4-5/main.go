// 写一个函数在原地完成消除[]string中相邻重复的字符串的操作。
package main

import "fmt"

func main() {
	arr := [...]string{"a", "b", "b", "b", "c", "c", "d"}
	sli := arr[:]
	for i := 1; i < len(sli); i++ {
		j := i
		for ; j < len(sli); j++ {
			if sli[j] != sli[j - 1] {
				break
			}
		}
		if i != j {
			copy(sli[i:], sli[j:])
			sli = sli[:len(sli) - j + i]
		}
	}
	fmt.Println(sli)
}
