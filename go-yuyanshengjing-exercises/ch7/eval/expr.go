package eval

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// numParams 参数数量
var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}

// Env 变量映射
type Env map[Var]float64

// Expr 计算
type Expr interface {
	// Eval returns the value of this Expr in the environment env.
	Eval(env Env) float64
	// Check reports errors in this Expr and adds its Vars to the set.
	Check(vars map[Var]float64) error
	// String 生成语法树
	String(env Env) string
}

// Var 识别一个变量
type Var string

// literal 字面量
type literal float64

// unary 一元运算
type unary struct {
	op rune // one of '+', '-'
	x  Expr
}

// binary 二元运算
type binary struct {
	op   rune // one of '+', '-', '*', '/'
	x, y Expr
}

// call 函数调用
type call struct {
	fn   string // one of "pow", "sin", "sqrt"
	args []Expr
}

// Eval 赋值运算
func (v Var) Eval(env Env) float64 {
	return env[v]
}

// Check 校验
func (v Var) Check(vars map[Var]float64) error {
	vars[v] = 0
	return nil
}

// String 语法树
func (v Var) String(env Env) string {
	return strconv.FormatFloat(env[v], 'g', -1, 64)
}

// Eval 字面量转换
func (l literal) Eval(env Env) float64 {
	return float64(l)
}

// Check 校验
func (literal) Check(vars map[Var]float64) error {
	return nil
}

// String 语法树
func (l literal) String(env Env) string {
	return strconv.FormatFloat(float64(l), 'g', -1, 64)
}

// Eval 一元运算
func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
}

// String 语法树
func (u unary) String(env Env) string {
	return fmt.Sprintf("%s%s", string(u.op), u.x.String(env))
}

// Check 校验
func (u unary) Check(vars map[Var]float64) error {
	if !strings.ContainsRune("+-", u.op) {
		return fmt.Errorf("unexpected unary op %q", u.op)
	}
	return u.x.Check(vars)
}

// Eval 一元运算
func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported binary operator: %q", b.op))
}

// Check 校验
func (b binary) Check(vars map[Var]float64) error {
	if !strings.ContainsRune("+-*/", b.op) {
		return fmt.Errorf("unexpected binary op %q", b.op)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

// String 语法树
func (b binary) String(env Env) string {
	return fmt.Sprintf("(%s %s %s)", b.x.String(env), string(b.op), b.y.String(env))
}

// Eval 函数调用
func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	case "max":
		if len(c.args) == 1 {
			return c.args[0].Eval(env)
		}
		c.args = c.args[1:]
		return math.Max(c.args[0].Eval(env), c.Eval(env))
	}
	panic(fmt.Sprintf("unsupported call operator: %s", c.fn))
}

// Check 校验
func (c call) Check(vars map[Var]float64) error {
	arity, ok := numParams[c.fn]
	if !ok {
		return fmt.Errorf("unknown function %q", c.fn)
	}
	if len(c.args) != arity {
		return fmt.Errorf("call to %s has %d args, want %d",
			c.fn, len(c.args), arity)
	}
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

// String 语法树
func (c call) String(env Env) string {
	args := make([]string, len(c.args))
	for i := 0; i < len(c.args); i++ {
		args[i] = c.args[i].String(env)
	}
	return fmt.Sprintf("%s("+strings.Join(args, ", ")+")", c.fn)
}

// ParseAndCheck 格式化并且校验
func ParseAndCheck(s string) (Expr, map[Var]float64, error) {
	if s == "" {
		return nil, nil, fmt.Errorf("empty expression")
	}
	expr, err := Parse(s)
	if err != nil {
		return nil, nil, err
	}
	vars := make(map[Var]float64)
	if err := expr.Check(vars); err != nil {
		return nil, nil, err
	}
	return expr, vars, nil
}
