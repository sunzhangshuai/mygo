package ch5

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
)

var sLinks []string
var labelCount map[string]int
var texts []string
var visitData = map[string]string{"a": "href", "img": "src", "link": "href", "script": "href"}

// visit 访问节点
func visit(n *html.Node) {
	if n == nil {
		return
	}

	if n.Type == html.ElementNode {
		labelCount[n.Data] += 1
		if key, ok := visitData[n.Data]; ok {
			for _, a := range n.Attr {
				if a.Key == key {
					sLinks = append(sLinks, a.Val)
				}
			}
		}
	} else if n.Type == html.TextNode {
		if n.Data != "script" && n.Data != "style" {
			texts = append(texts, n.Data)
		}
	}

	visit(n.FirstChild)
	visit(n.NextSibling)
	return
}

// printData 打印数据
func printData() {
	lks, _ := json.MarshalIndent(sLinks, "", "\t")
	fmt.Println("sLinks：", string(lks))
	lc, _ := json.MarshalIndent(labelCount, "", "\t")
	fmt.Println("labelCount：", string(lc))
	ltx, _ := json.MarshalIndent(texts, "", "\t")
	fmt.Println("texts：", string(ltx))
}

// init 初始化
func init() {
	sLinks = make([]string, 0, 100)
	labelCount = make(map[string]int)
	texts = make([]string, 0, 100)
}
