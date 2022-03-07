// 编写一个函数，计算两个SHA256哈希码中不同bit的数目。
package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

func main() {
	c1 := sha256.Sum256([]byte("张三"))
	c2 := sha256.Sum256([]byte("李四"))
	fmt.Println(c1)
	fmt.Println(c2)


	ret := 0
	for i, val := range c1 {
		if val != c2[i] {
			ret++
			fmt.Println("==============="+ strconv.Itoa(ret) +"================")

			fmt.Println(val)
			fmt.Println(c2[i])
			fmt.Println("===============================")

		}
	}
	fmt.Println(len(c1))
	fmt.Println(ret)
	// Output:
	// 2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881
	// 4b68ab3847feda7d6c62c1fbcbeebfa35eab7351ed5e78f4ddadea5df64b8015
	// false
	// [32]uint8
}
