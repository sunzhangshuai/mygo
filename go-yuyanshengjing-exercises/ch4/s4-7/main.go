// 修改reverse函数用于原地反转UTF-8编码的[]byte。是否可以不用分配额外的内存？
package main

import (
	"fmt"
)

// 用的非原地
func reverse(s []byte) []byte {
	re := make([]byte, len(s), len(s))
	start := 0
	for sLen := len(s); sLen > 0; {
		size := 0
		if s[start] < 192 { // 0开头
			size = 1
		} else if s[start] < 224 { // 110 开头
			size = 2
		}else if s[start] < 240 { // 1110 开头
			size = 3
		} else {
			size = 4
		}
		copy(re[sLen - size: sLen], s[start:start + size])
		start += size
		sLen -= size
	}
	return re
}

func main()  {
	str := "丑ـبــثわたし"
	s := []byte(str)
	fmt.Println(s)
	fmt.Println(string(reverse(s)))
}