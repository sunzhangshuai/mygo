// 编写一个程序，默认情况下打印标准输入的SHA256编码，并支持通过命令行flag定制，输出SHA384或SHA512哈希算法。
package main

import (
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"
)

var(
	h string
	s string
)

func init() {
	flag.StringVar(&s, "s", "", "设置要编码的字符串")
	flag.StringVar(&h, "h", "", "设置编码格式，option：sha384|sha512")
}

func main() {
	flag.Parse()
	if h != "sha384" && h != "sha512" || len(s) == 0 {
		flag.Usage()
		return
	}

	if h == "sha384" {
		sha := sha512.Sum384([]byte(s))
		sh := sha[:]
		fmt.Println(hex.EncodeToString(sh))
	} else {
		sha := sha512.Sum512([]byte(s))
		sh := sha[:]
		fmt.Println(hex.EncodeToString(sh))
	}
}
