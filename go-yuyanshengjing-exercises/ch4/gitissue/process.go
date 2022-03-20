package gitissue

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/go-github/v43/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"os"
	"strconv"
	"strings"
)

// Repository 项目名称
const Repository = "mygo"

// Owner 作者
const Owner = "sunzhangshuai"

// AccessToken token
const AccessToken = "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"

// client github 客户端
var client *github.Client

// Run git issue 处理
func Run() error {
	ctx := context.Background()

	initClient(ctx)

	var err error

	scanner := bufio.NewScanner(os.Stdin)

	// 主站提示
	if err = homeTip(); err != nil {
		return err
	}

	// 始终读取数据
	for scanner.Scan() {
		action := strings.TrimSpace(scanner.Text())

		switch action {
		case "0":
			return nil
		case "1":
			err = create(ctx, scanner)
		case "2":
			err = list(ctx, scanner)
		case "3":
			err = update(ctx, scanner)
		case "4":
			err = deleteIssue(ctx, scanner)
		default:
			return errors.New("input error")
		}
		if err != nil {
			return err
		}
		if err = homeTip(); err != nil {
			return err
		}
	}
	return nil
}

// create 创建issue
func create(ctx context.Context, scanner *bufio.Scanner) error {
	var err error

	fmt.Println("请输入标题")
	scanner.Scan()
	title := scanner.Text()
	fmt.Println("请输入内容")
	scanner.Scan()
	body := scanner.Text()

	request := &github.IssueRequest{
		Title: &title,
		Body:  &body,
	}

	issue, r, err := client.Issues.Create(ctx, Owner, Repository, request)
	if err != nil {
		return errors.Wrap(err, "create err")
	}
	defer r.Body.Close()
	if err = issueTip(issue); err != nil {
		return errors.Wrap(err, "create err")
	}
	return nil
}

// list 读取issue列表
func list(ctx context.Context, scanner *bufio.Scanner) error {
	fun := func(id int) error {
		// 处理
		issue, itemRes, err := client.Issues.Get(ctx, Owner, Repository, id)
		if err != nil {
			return errors.Wrap(err, "cat item err")
		}
		defer itemRes.Body.Close()

		if err = issueTip(issue); err != nil {
			return errors.Wrap(err, "list err")
		}
		return nil
	}
	return processIssue(ctx, scanner, fun)
}

// deleteIssue 删除issue
func deleteIssue(ctx context.Context, scanner *bufio.Scanner) error {
	fun := func(id int) error {
		// 处理
		clos := "close"
		issue, itemRes, err := client.Issues.Edit(ctx, Owner, Repository, id, &github.IssueRequest{State: &clos})
		if err != nil {
			return errors.Wrap(err, "cat item err")
		}
		defer itemRes.Body.Close()
		fmt.Println(123213, issue)
		if err = issueTip(issue); err != nil {
			return errors.Wrap(err, "list err")
		}
		return nil
	}
	return processIssue(ctx, scanner, fun)
}

// update 修改issue
func update(ctx context.Context, scanner *bufio.Scanner) error {
	fun := func(id int) error {
		fmt.Println("请输入标题")
		scanner.Scan()
		title := scanner.Text()
		fmt.Println("请输入内容")
		scanner.Scan()
		body := scanner.Text()
		request := &github.IssueRequest{
			Title: &title,
			Body:  &body,
		}
		issue, itemRes, err := client.Issues.Edit(ctx, Owner, Repository, id, request)
		if err != nil {
			return errors.Wrap(err, "cat item err")
		}
		defer itemRes.Body.Close()
		if err = issueTip(issue); err != nil {
			return errors.Wrap(err, "list err")
		}
		return nil
	}
	return processIssue(ctx, scanner, fun)
}

// processItem 处理单个issue
func processIssue(ctx context.Context, scanner *bufio.Scanner, f func(id int) error) error {
	// 1.获取list
	issues, r, err := client.Issues.ListByRepo(ctx, Owner, Repository, nil)
	if err != nil {
		return errors.Wrap(err, "list err")
	}
	defer r.Body.Close()

	// 2.输出list
	if err = issueListTip(issues); err != nil {
		return errors.Wrap(err, "list err")
	}

	// 3. 没有数据要自动退出
	if len(issues) == 0 {
		return nil
	}

	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	// 退出
	if input == "0" {
		return nil
	}

	var id int
	flag := false
	// 验证输入
	if id, err = strconv.Atoi(input); err != nil {
		return errors.Wrap(err, "list input err")
	}
	for _, item := range issues {
		if *item.Number == id {
			flag = true
		}
	}
	if !flag {
		return errors.New("list input err")
	}

	return f(id)
}

// 初始化客户端
func initClient(ctx context.Context) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
}
