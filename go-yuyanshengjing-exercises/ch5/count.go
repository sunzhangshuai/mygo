package ch5

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

// CountWordsAndImages 获取单词和图片个数
func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}
	words, images = countWordsAndImages(doc)
	return
}

// cuntWordsAndImages 获取单词和图片个数
func countWordsAndImages(n *html.Node) (words, images int) {
	if n == nil {
		return
	}
	if n.Type == html.ElementNode && n.Data == "img" {
		for _, a := range n.Attr {
			if a.Key == "src" {
				images++
			}
		}
	} else if n.Type == html.TextNode {
		if n.Data != "script" && n.Data != "style" {
			reader := strings.NewReader(n.Data)

			buf := bufio.NewScanner(reader)
			buf.Split(bufio.ScanWords)
			for buf.Scan() {
				words++
			}
		}
	}
	words1, images1 := countWordsAndImages(n.FirstChild)
	words2, images2 := countWordsAndImages(n.NextSibling)
	words += words1 + words2
	images += images1 + images2
	return
}
