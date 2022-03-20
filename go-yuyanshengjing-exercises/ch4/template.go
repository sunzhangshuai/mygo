package ch4

import (
	"context"
	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
	"html/template"
	"io"
)

// AccessToken token
const AccessToken = "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"

// IssueTemp html模板
const IssueTemp = `
<h1>{{.Total}} issues</h1>
<table>
<tr style='text-align: left'>
  <th>#</th>
  <th>State</th>
  <th>User</th>
  <th>Title</th>
</tr>
{{range .Issues}}
<tr>
  <td><a href='{{.HTMLURL}}'>{{.Number}}</a></td>
  <td>{{.State}}</td>
  <td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
  <td><a href='{{.HTMLURL}}'>{{.Title}}</a></td>
</tr>
{{end}}
</table>
`

// client github客户端
var client *github.Client

// RunTemplate 处理模板
func RunTemplate(writer io.Writer) error {
	request := &github.SearchOptions{}
	request.PerPage = 1000
	issues, r, err := client.Search.Issues(context.Background(), "repo:golang/go json decoder", request)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	issuesTem, err := template.New("issues").Parse(IssueTemp)
	if err != nil {
		return err
	}
	return issuesTem.Execute(writer, issues)
}

// init 初始化客户端
func init() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: AccessToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client = github.NewClient(tc)
}
