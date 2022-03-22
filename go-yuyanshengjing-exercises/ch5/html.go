package ch5

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

// forEachNode针对每个结点x，都会调用pre(x)和post(x)。
// pre和post都是可选的。
// 遍历孩子结点之前，pre被调用
// 遍历孩子结点之后，post被调用
func forEachNode(n *html.Node, pre, post func(n *html.Node, f func(ID string) bool, depth func(action string) int) bool, f func(ID string) bool, depth func(action string) int) bool {
	if pre != nil {
		if !pre(n, f, depth) {
			return false
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if !forEachNode(c, pre, post, f, depth) {
			return false
		}
	}
	if post != nil {
		return post(n, f, depth)
	}
	return true
}

func startElement(n *html.Node, f func(ID string) bool, depth func(action string) int) bool {
	if n.Type == html.ElementNode {
		attrs := make([]string, len(n.Attr))
		for i, attr := range n.Attr {
			if attr.Key == "id" && f(attr.Val) {
				return false
			}
			attrs[i] = " " + attr.Key + "=\"" + attr.Val + "\""
		}
		attrStr := strings.Join(attrs, "")
		if n.FirstChild == nil {
			fmt.Printf("%*s<%s%s/>\n", depth("add"), "", n.Data, attrStr)
		} else {
			fmt.Printf("%*s<%s%s>\n", depth("add"), "", n.Data, attrStr)
		}
	}
	if n.Type == html.TextNode {
		fmt.Printf("%s\n", n.Data)
	}
	return true
}

func endElement(n *html.Node, f func(ID string) bool, depth func(action string) int) bool {
	if n.Type == html.ElementNode && n.FirstChild != nil {
		fmt.Printf("%*s</%s>\n", depth("sub"), "", n.Data)
	}
	return true
}
