package mirror

import (
	"crypto/md5"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	http "net/http"
	"net/url"
	"os"
	"path"
	"runtime"
)

// URLInfo 格式
type URLInfo struct {
	depth int
	data  string
	Type  string
}

var basePath string

// init 初始化
func init() {
	_, fileName, _, _ := runtime.Caller(0)
	basePath = path.Join(path.Dir(fileName), "file")
}

// Run 镜像处理
func Run(depth int, reqURL string) error {
	urlInfo := &URLInfo{0, reqURL, "html"}

	workList := make(chan []*URLInfo)  // URL列表，可能有重复项
	unseenLinks := make(chan *URLInfo) // 消除重复的URL

	// 将命令行参数添加到工作列表
	go func() { workList <- []*URLInfo{urlInfo} }()

	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				newList, err := crawl(link)
				if err != nil {
					log.Println(err)
				}
				workList <- newList
			}
		}()
	}

	seen := make(map[string]bool)
	for list := range workList {
		for _, link := range list {
			if !seen[link.data] && link.depth < depth {
				seen[link.data] = true
				unseenLinks <- link
			}
		}
	}
	return nil
}

// crawl 爬虫
func crawl(urlInfo *URLInfo) ([]*URLInfo, error) {
	var body io.ReadCloser
	var err error
	var urlList []*URLInfo
	var doc *html.Node

	// 拿到数据
	if body, err = request(urlInfo.data); err != nil {
		return nil, err
	}
	defer body.Close()

	// 获取url列表
	if urlInfo.Type == "html" {
		if doc, urlList, err = extract(body, urlInfo); err != nil {
			return nil, err
		}
	}

	// 保存文件
	if err = saveFile(makeLocalUrl(urlInfo), doc, body); err != nil {
		return nil, err
	}
	return urlList, nil
}

// makeRemoteUrl 获取远程url
func makeRemoteUrl(parentURL *URLInfo, reqURL string, reqType string) (*URLInfo, error) {
	var result, parentInfo *url.URL
	var err error

	if result, err = url.Parse(reqURL); err != nil {
		return nil, err
	}
	if result.Host == "" {
		if parentInfo, err = url.Parse(parentURL.data); err != nil {
			return nil, err
		}
		parentInfo.Path = path.Join(path.Dir(parentInfo.Path), reqURL)
		return &URLInfo{parentURL.depth + 1, parentInfo.String(), reqType}, nil
	}
	return &URLInfo{parentURL.depth + 1, result.String(), reqType}, nil
}

// makeLocalUrl 获取本地url
func makeLocalUrl(urlInfo *URLInfo) *URLInfo {
	urlData := path.Join(basePath, urlInfo.Type, fmt.Sprintf("%x", md5.Sum([]byte(urlInfo.data)))+"."+urlInfo.Type)
	return &URLInfo{urlInfo.depth, urlData, urlInfo.Type}
}

// request 请求
func request(reqURL string) (io.ReadCloser, error) {
	var resp *http.Response
	var err error

	if resp, err = http.Get(reqURL); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("request fail")
	}
	return resp.Body, nil
}

// saveFile 保存文件
func saveFile(urlInfo *URLInfo, doc *html.Node, body io.ReadCloser) error {
	var file *os.File
	var err error

	if err = CheckDir(path.Dir(urlInfo.data)); err != nil {
		return err
	}
	if file, err = os.Create(urlInfo.data); err != nil {
		return err
	}
	defer file.Close()

	// 保存文件
	if urlInfo.Type == "html" {
		html.Render(file, doc)
	} else {
		io.Copy(file, body)
	}

	return nil
}

// CheckDir 校验目录
func CheckDir(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else {
		err := os.MkdirAll(path, 0711)
		if err != nil {
			return err
		}
	}
	// check again
	_, err := os.Stat(path)
	return err
}

// extract makes an HTTP GET request to the specified URL, parses
// the response as HTML, and returns the links in the HTML document.
func extract(body io.ReadCloser, urlInfo *URLInfo) (*html.Node, []*URLInfo, error) {
	var doc *html.Node
	var err error

	// 格式化输入
	if doc, err = html.Parse(body); err != nil {
		return nil, nil, fmt.Errorf("parsing %s as HTML: %v", urlInfo.data, err)
	}

	// 循环处理
	list := make([]*URLInfo, 0)
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "img") {
			for k, a := range n.Attr {
				switch a.Key {
				case "href":
					if path.Ext(a.Val) == ".pdf" {
						appendURL(urlInfo, &a, path.Ext(a.Val)[1:], &list)
					} else {
						appendURL(urlInfo, &a, "html", &list)
					}
				case "src":
					if path.Ext(a.Val) != "" {
						appendURL(urlInfo, &a, path.Ext(a.Val)[1:], &list)
					}
				}
				n.Attr[k] = a
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return doc, list, nil
}

// appendURL 追加url
func appendURL(parentUrlInfo *URLInfo, attribute *html.Attribute, urlType string, list *[]*URLInfo) {
	var remoteURL *URLInfo
	var err error

	if remoteURL, err = makeRemoteUrl(parentUrlInfo, attribute.Val, urlType); err != nil {
		return
	}

	// 赋值为本地url
	attribute.Val = makeLocalUrl(remoteURL).data
	*list = append(*list, remoteURL)
}

// Copied from gopl.io/ch5/outline2.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
