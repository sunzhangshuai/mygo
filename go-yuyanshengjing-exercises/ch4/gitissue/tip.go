package gitissue

import (
	"fmt"
	"github.com/google/go-github/v43/github"
	"github.com/pkg/errors"
	"os"
	"text/template"
	"time"
)

const (
	// HomeTemp 首页模板
	HomeTemp = `请选择操作类型，操作账户【{{.}}】：
	1. 创建issue
	2. 读取issue
	3. 更新issue
	4. 关闭issue
`

	// ListTemp 列表模板
	ListTemp = `{{.Total}} issues:
----------------------------------------
{{range .Items}}{{.Number}}. {{.Title | showTitle }} {{.User.Login}} {{.CreatedAt | daysAgo}} days
{{end}}----------------------------------------
{{if eq .Total 0 }}没有issue，自动返回上一级。
{{else}} 请选择要处理的issue，输入0退出。
{{end}}`

	// IssueTemp 议题模板
	IssueTemp = `序号：{{.Number}}. 
标题：{{.Title | showTitle }} 
用户：{{.User.Login}} 
创建时间：{{.CreatedAt | daysAgo}} days
地址：{{.URL}} 
HTML地址：{{.HTMLURL}} 
状态：{{.State}} 
`
)

type ListData struct {
	Total int
	Items []*github.Issue
}

// daysAgo 多少天前
func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

func showTitle(s *string) string {
	return fmt.Sprintf("%.64s", *s)
}

// homeTip 首页提示
func homeTip() error {
	var tip *template.Template
	var err error
	if tip, err = template.New("repo").Parse(HomeTemp); err != nil {
		return errors.Wrap(err, "home err")
	}
	if err = tip.Execute(os.Stdout, Repository); err != nil {
		return errors.Wrap(err, "home err")
	}
	return nil
}

// issueListTip issue 列表提示
func issueListTip(issues []*github.Issue) error {
	var tip *template.Template
	var err error

	dataTemp := &ListData{
		Total: len(issues),
		Items: issues,
	}
	if tip, err = template.New("repo").Funcs(template.FuncMap{"daysAgo": daysAgo, "showTitle": showTitle}).Parse(ListTemp); err != nil {
		return errors.Wrap(err, "list err")
	}
	if err = tip.Execute(os.Stdout, dataTemp); err != nil {
		return errors.Wrap(err, "list err")
	}
	return nil
}

// issueTip issue 数据
func issueTip(issue *github.Issue) error {
	var tip *template.Template
	var err error
	if tip, err = template.New("repo").Funcs(template.FuncMap{"daysAgo": daysAgo, "showTitle": showTitle}).Parse(IssueTemp); err != nil {
		return errors.Wrap(err, "item err")
	}
	if err = tip.Execute(os.Stdout, issue); err != nil {
		return errors.Wrap(err, "list err")
	}
	return nil
}
