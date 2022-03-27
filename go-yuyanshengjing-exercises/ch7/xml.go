package ch7

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

// XMLSelect xml 选择器
func XMLSelect(r io.Reader, attrKey, attrValue string) error {
	dec := xml.NewDecoder(r)
	var stack []string // stack of element names
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			stack = append(stack, tok.Name.Local) // push
			for _, attr := range tok.Attr {
				if attr.Name.Local == attrKey && attr.Value == attrValue {
					fmt.Println(stack)
				}
			}
		case xml.EndElement:
			stack = stack[:len(stack)-1] // pop
		}
	}
	return nil
}

// ParseTree 格式化树
func ParseTree(r io.Reader) (*Element, error) {
	dec := xml.NewDecoder(r)
	var stack []*Element
	stack = append(stack, &Element{}) // push
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			elem := &Element{
				Type: tok.Name,
				Attr: tok.Attr,
			}
			stack = append(stack, elem) // push
		case xml.EndElement:
			stack[len(stack)-2].Children = append(stack[len(stack)-2].Children, stack[len(stack)-1])
			stack = stack[:len(stack)-1] // pop
		case xml.CharData:
			str := strings.TrimSpace(string(tok))
			if str != "" {
				stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, str)
			}
		}
	}
	return stack[0].Children[0].(*Element), nil
}

type Node interface{} // CharData or *Element

type CharData string

type Element struct {
	Type     xml.Name
	Attr     []xml.Attr
	Children []Node
}

// containsAll reports whether x contains the elements of y, in order.
func containsAll(x, y []string) bool {
	for len(y) <= len(x) {
		if len(y) == 0 {
			return true
		}
		if x[0] == y[0] {
			y = y[1:]
		}
		x = x[1:]
	}
	return false
}
