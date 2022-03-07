// 编写一个非递归版本的comma函数，使用bytes.Buffer代替字符串链接操作。
package main

import (
	"bytes"
	"fmt"
)

func main() {
	fmt.Println(comma("12321321321"))
}

func comma(s string) string {
	var buf bytes.Buffer
	strLength := len(s)
	for i := strLength; i > 0; i-- {
		buf.WriteString(s[strLength-i:strLength-i+1])
		if i != 1 && (i - 1) % 3 == 0 {
			buf.WriteString(",")
		}
	}
	return buf.String()
}
