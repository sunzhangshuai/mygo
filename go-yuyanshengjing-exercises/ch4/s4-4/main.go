// 编写一个rotate函数，通过一次循环完成旋转。
package main

import "fmt"

func rotate(s []int, l int)  {
	zero := s[0]
	sLen := len(s)
	for i := 0; i < sLen; i = i + 1 {
		fmt.Println((i - l + sLen) % sLen)
		s[i] = s[(i - l + sLen) % sLen]
	}
	s[(l + sLen) % sLen] = zero
}

func main()  {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	rotate(s, 1)
	fmt.Println(s)
}
