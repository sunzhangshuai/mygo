package ch5

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

// Exercises 练习
type Exercises struct {
}

// findLinks 获取链接
func findLinks() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}
	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}
}

// visit appends to links each link found in n and returns the result.
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}

func countWordsAndImages(n *html.Node) (words, images int) { /* ... */ }

// Task1 修改findlinks代码中遍历n.FirstChild链表的部分，将循环调用visit，改成递归调用。
func (e *Exercises) Task1() {
	for _, link := range e.task1Visit(nil, e.getDoc()) {
		fmt.Println(link)
	}
}

// Task2 编写函数，记录在HTML树中出现的同名元素的次数。
func (e *Exercises) Task2() {
	links := make(map[string]int)
	for label, count := range e.task2Visit(links, e.getDoc()) {
		fmt.Println(label, count)
	}
}

// Task3 编写函数输出所有text结点的内容。注意不要访问<script>和<style>元素，因为这些元素对浏览者是不可见的。
func (e *Exercises) Task3() {
	for _, text := range e.task3Visit(nil, e.getDoc()) {
		fmt.Println(text)
	}
}

// Task4 扩展visit函数，使其能够处理其他类型的结点，如images、scripts和style sheets。
func (e *Exercises) Task4() {
	for _, link := range e.task4Visit(nil, e.getDoc()) {
		fmt.Println(link)
	}
}

// Task5 实现countWordsAndImages。
func (e *Exercises) Task5() {
	if len(os.Args) < 6 {
		fmt.Println("需要url参数")
		return
	}
	url := os.Args[5]
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

// getDoc 获取文档
func (e *Exercises) getDoc() *html.Node {
	_, fullFilename, _, _ := runtime.Caller(0)
	path := filepath.Dir(fullFilename)
	file, err := os.OpenFile(path+"/input/text.html", os.O_RDWR, 0777)
	doc, err := html.Parse(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks: %v\n", err)
		os.Exit(1)
	}
	return doc
}

// task1Visit 修改findlinks代码中遍历n.FirstChild链表的部分，将循环调用visit，改成递归调用。
func (e *Exercises) task1Visit(links []string, n *html.Node) []string {
	if n == nil {
		return links
	}
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	links = e.task1Visit(links, n.FirstChild)
	links = e.task1Visit(links, n.NextSibling)
	return links
}

// task2Visit 编写函数，记录在HTML树中出现的同名元素的次数。
func (e *Exercises) task2Visit(links map[string]int, n *html.Node) map[string]int {
	if n == nil {
		return links
	}
	if n.Type == html.ElementNode {
		links[n.Data] += 1
	}
	links = e.task2Visit(links, n.FirstChild)
	links = e.task2Visit(links, n.NextSibling)
	return links
}

// task3Visit 编写函数输出所有text结点的内容。注意不要访问<script>和<style>元素，因为这些元素对浏览者是不可见的。
func (e *Exercises) task3Visit(texts []string, n *html.Node) []string {
	if n == nil || n.Data == "script" || n.Data == "style" {
		return texts
	}
	if n.Type == html.TextNode {
		texts = append(texts, n.Data)
	}
	texts = e.task3Visit(texts, n.FirstChild)
	texts = e.task3Visit(texts, n.NextSibling)
	return texts
}

var visitData = map[string]bool{"a": true, "img": true, "link": true, "script": true}

// task4Visit 扩展visit函数，使其能够处理其他类型的结点，如images、scripts和style sheets。
func (e *Exercises) task4Visit(links []string, n *html.Node) []string {
	if n == nil {
		return links
	}
	if n.Type == html.ElementNode && visitData[n.Data] {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	links = e.task4Visit(links, n.FirstChild)
	links = e.task4Visit(links, n.NextSibling)
	return links
}
