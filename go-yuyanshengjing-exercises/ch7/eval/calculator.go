package eval

import (
	"html/template"
	"io"
	"path/filepath"

	"exercises/util"
)

// Run 执行结果
func Run(writer io.Writer, expr string) error {
	type Data struct {
		Expr   string
		Result float64
	}

	data := Data{
		Expr: expr,
	}

	tem, err := template.ParseFiles(util.GetFileName(1, filepath.Join("../file", "calculator.htm")))
	if err != nil {
		return err
	}

	if expr != "" {
		parse, err := Parse(expr)
		if err != nil {
			return err
		}
		ret := parse.Eval(nil)

		data.Result = ret
	}
	return tem.Execute(writer, data)
}
