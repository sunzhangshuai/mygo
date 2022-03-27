package ch7

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"

	"exercises/util"
)

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database struct {
	data sync.Map
}

var DB database

// List 数据库列表
func (db *database) List(w http.ResponseWriter, req *http.Request) {
	db.list(w)
}

func (db *database) list(w io.Writer) {
	data := make(map[string]dollars)
	temp, err := template.ParseFiles(util.GetFileName(1, filepath.Join("file", "database.htm")))
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	db.data.Range(func(key, value interface{}) bool {
		data[key.(string)] = dollars(value.(float32))
		return true
	})

	if err = temp.Execute(w, data); err != nil {
		fmt.Fprintln(w, err)
		return
	}
}

// Add db增加数据
func (db *database) Add(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")
	fPrice, err := strconv.ParseFloat(price, 32)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	db.data.Store(item, float32(fPrice))
}

// Delete db删除数据
func (db *database) Delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	db.data.Delete(item)
}

// Update db修改数据
func (db *database) Update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")
	fPrice, err := strconv.ParseFloat(price, 32)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	if _, ok := db.data.Load(item); !ok {
		fmt.Fprintln(w, item+"不存在")
		return
	}
	db.data.Store(item, float32(fPrice))
}

// Get 获取价格
func (db *database) Get(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if price, ok := db.data.Load(item); !ok {
		fmt.Fprintln(w, item+"不存在")
		return
	} else {
		fmt.Fprintf(w, "%s\n", dollars(price.(float32)))
	}
}

// init 初始化数据库
func init() {
	DB = database{}
	DB.data.Store("shoes", float32(50))
	DB.data.Store("socks", float32(5))
}
