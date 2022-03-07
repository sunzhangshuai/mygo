// 修改issues程序，根据问题的时间进行分类，比如不到一个月的、不到一年的、超过一年。
package main

import (
	"fmt"
	"gopl.io/ch4/github"
	"log"
	"time"
)

func main() {
	args := []string{"repo:golang/go", "is:open", "json decoder"}
	result, err := github.SearchIssues(args)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	res := make(map[string][]string)
	res["far"] = make([]string, 0)
	res["middle"] = make([]string, 0)
	res["near"] = make([]string, 0)
	nowTime := time.Now()
	for _, item := range result.Items {
		createdAt := item.CreatedAt
		if createdAt.AddDate(1, 0, 0).Before(nowTime)  {
			res["far"] = append(res["far"], item.HTMLURL)
		} else if createdAt.AddDate(0, 1, 0).Before(nowTime) {
			res["middle"] = append(res["far"], item.HTMLURL)
		} else {
			res["near"] = append(res["far"], item.HTMLURL)
		}
	}
	fmt.Println(res)
}