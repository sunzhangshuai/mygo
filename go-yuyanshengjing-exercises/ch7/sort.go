package ch7

import (
	"html/template"
	"io"
	"path"
	"sort"
	"time"

	"exercises/util"
)

type byArtist []*Track

var order1 string
var order2 string
var orderFunc = make(map[string]func(x byArtist, i, j int) int, 5)

// Track 音乐
type Track struct {
	Title  string        // 标题
	Artist string        // 演唱者
	Album  string        // 专辑
	Year   int           // 出版时间
	Length time.Duration // 音乐长度
}

// tracks 播放列表
var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

// length 歌曲长度
func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

// Len 长度
func (x byArtist) Len() int { return len(x) }

// Less 顺序
func (x byArtist) Less(i, j int) bool {
	if order1 != "" {
		if orderFunc[order1](x, i, j) == -1 {
			return true
		} else if orderFunc[order1](x, i, j) == 1 {
			return false
		}
	}
	if order2 != "" {
		if orderFunc[order1](x, i, j) == -1 {
			return true
		} else {
			return false
		}
	}
	return true
}

// Swap 交换方式
func (x byArtist) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

// init 初始化
func init() {
	orderFunc["Title"] = func(x byArtist, i, j int) int {
		if x[i].Title > x[j].Title {
			return 1
		} else if x[i].Title < x[j].Title {
			return -1
		}
		return 0
	}
	orderFunc["Artist"] = func(x byArtist, i, j int) int {
		if x[i].Artist > x[j].Artist {
			return 1
		} else if x[i].Artist < x[j].Artist {
			return -1
		}
		return 0
	}
	orderFunc["Album"] = func(x byArtist, i, j int) int {
		if x[i].Album > x[j].Album {
			return 1
		} else if x[i].Album < x[j].Album {
			return -1
		}
		return 0
	}
	orderFunc["Year"] = func(x byArtist, i, j int) int {
		if x[i].Year > x[j].Year {
			return 1
		} else if x[i].Year < x[j].Year {
			return -1
		}
		return 0
	}
	orderFunc["Length"] = func(x byArtist, i, j int) int {
		if x[i].Length > x[j].Length {
			return 1
		} else if x[i].Length < x[j].Length {
			return -1
		}
		return 0
	}
}

// FmtTracks 格式化播放列表
func FmtTracks(writer io.Writer, order string) error {
	order2, order1 = order1, order
	temp, err := template.ParseFiles(util.GetFileName(1, path.Join("file", "track.htm")))
	if err != nil {
		return err
	}
	sort.Stable(byArtist(tracks))
	return temp.Execute(writer, tracks)
}

// IsPalindrome 判断是否为回文序列
func IsPalindrome(s sort.Interface) bool {
	i := 0
	j := s.Len() - 1
	for i < j {
		if s.Less(i, j) || s.Less(j, i) {
			return false
		}
		i++
		j--
	}
	return true
}
