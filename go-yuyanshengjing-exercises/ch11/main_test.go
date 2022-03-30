package ch11

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"exercises/ch6"
)

// TestTask1 为4.3节中的charcount程序编写测试。
func TestTask1(t *testing.T) {
	var tests = []struct {
		input string
		want  map[rune]int
	}{
		{"aaavvvccc", map[rune]int{'a': 3, 'v': 3, 'c': 3}},
		{"aaa张ccc", map[rune]int{'a': 3, '张': 1, 'c': 3}},
	}
	for _, test := range tests {
		if got := charCount(strings.NewReader(test.input)); !reflect.DeepEqual(got, test.want) {
			t.Errorf("charCount(%s) = %v", test.input, got)
		}
	}
}

// TestTask2 为（§6.5）的IntSet编写一组测试，用于检查每个操作后的行为和基于内置map的集合等价，后面练习11.7将会用到。
func TestTask2(t *testing.T) {
	var x, y IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)

	y.Add(9)
	y.Add(42)

	var x1, y1 ch6.IntSet
	x1.Add(1)
	x1.Add(144)
	x1.Add(9)

	y1.Add(9)
	y1.Add(42)

	var tests = []struct {
		input func() string
		want  func() string
	}{
		{func() string {
			x.IntersectWith(&y)
			return x.String()
		}, func() string {
			x1.IntersectWith(&y1)
			return x1.String()
		}},
		{func() string {
			x.UnionWith(&y)
			return x.String()
		}, func() string {
			x1.UnionWith(&y1)
			return x1.String()
		}},
		{func() string {
			x.DifferenceWith(&y)
			return x.String()
		}, func() string {
			x1.DifferenceWith(&y1)
			return x1.String()
		}},
		{func() string {
			x.SymmetricDifference(&y)
			return x.String()
		}, func() string {
			x1.SymmetricDifference(&y1)
			return x1.String()
		}},
	}
	for _, test := range tests {
		if test.input() != test.want() {
			t.Errorf("%s != %s", x.String(), x1.String())
		}
	}
}

// TestTask3 TestRandomPalindromes测试函数只测试了回文字符串。编写新的随机测试生成器，用于测试随机生成的非回文字符串。
func TestTask3(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 1000; i++ {
		p := randomPalindrome(rng)
		if !IsPalindrome(p) {
			t.Errorf("IsPalindrome(%q) = false", p)
		}
	}

	for i := 0; i < 1000; i++ {
		p := notRandomPalindrome(rng)
		if IsPalindrome(p) {
			t.Errorf("IsPalindrome(%q) = true", p)
		}
	}
}

// TestTask4 修改randomPalindrome函数，以探索IsPalindrome是否对标点和空格做了正确处理。
func TestTask4(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 1000; i++ {
		p := randomPalindrome(rng)
		if !IsPalindrome(p) {
			t.Errorf("IsPalindrome(%q) = false", p)
		}
	}
}

// TestTask5 用表格驱动的技术扩展TestSplit测试，并打印期望的输出结果。
func TestTask5(t *testing.T) {
	var tests = []struct {
		s    string
		sep  string
		want int
	}{
		{"a:b:c", ":", 3},
		{"aa:b:c", ":", 3},
		{"aa:b:c", "b", 2},
	}
	for _, test := range tests {
		words := strings.Split(test.s, test.sep)
		if got := len(words); got != test.want {
			t.Errorf("Split(%q, %q) returned %d words, want %d",
				test.s, test.sep, got, test.want)
		}
	}
}

// benchmarkTask6 为2.6.2节的练习2.4和练习2.5的PopCount函数编写基准测试。看看基于表格算法在不同情况下对提升性能会有多大帮助。
func benchmarkTask6(b *testing.B, typ string) {
	for i := 0; i >= b.N; i++ {
		switch typ {
		case "table":
			PopCountTable(99999)
		case "bit":
			PopCountBit(99999)
		case "bitAnd":
			PopCountBitAnd(99999)
		}
	}
}

func BenchmarkTask6Table(b *testing.B) {
	benchmarkTask6(b, "table")
}

func BenchmarkTask6Bit(b *testing.B) {
	benchmarkTask6(b, "bit")
}

func BenchmarkTask6BitAnd(b *testing.B) {
	benchmarkTask6(b, "bitAnd")
}

// benchmarkTask7 为*IntSet（§6.5）的Add、UnionWith和其他方法编写基准测试，使用大量随机输入。
// 你可以让这些方法跑多快？选择字的大小对于性能的影响如何？IntSet和基于内建map的实现相比有多快？
func BenchmarkTask7Bit(b *testing.B) {
	var x ch6.IntSet
	seed := time.Now().UTC().UnixNano()
	b.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i >= b.N; i++ {
		x.Add(rng.Intn(10000))
	}
}

func BenchmarkTask7Map(b *testing.B) {
	var x IntSet
	seed := time.Now().UTC().UnixNano()
	b.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i >= b.N; i++ {
		x.Add(rng.Intn(10000))
	}
}

func BenchmarkTask7BitUnion(b *testing.B) {
	seed := time.Now().UTC().UnixNano()
	b.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < b.N; i++ {
		var x, y ch6.IntSet
		for i := 0; i <= 100; i++ {
			x.Add(rng.Intn(10000))
			y.Add(rng.Intn(10000))
		}
		x.UnionWith(&y)
	}
}

func BenchmarkTask7MapUnion(b *testing.B) {
	seed := time.Now().UTC().UnixNano()
	b.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < b.N; i++ {
		var x, y IntSet
		for i := 0; i <= 100; i++ {
			x.Add(rng.Intn(10000))
			y.Add(rng.Intn(10000))
		}
		x.UnionWith(&y)
	}
}
