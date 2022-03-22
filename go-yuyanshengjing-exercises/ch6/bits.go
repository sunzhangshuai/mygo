package ch6

import (
	"bytes"
	"fmt"

	"exercises/ch2"
)

// IntSet bit 数组
type IntSet struct {
	words []uint
}

const size = 32 << (^uint(0) >> 63)

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/size, uint(x%size)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/size, uint(x%size)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tWord := range t.words {
		if i < len(s.words) {
			s.words[i] |= tWord
		} else {
			s.words = append(s.words, tWord)
		}
	}
}

// IntersectWith 交集：元素在A集合B集合均出现
func (s *IntSet) IntersectWith(t *IntSet) {
	var sLen int
	for i, tWord := range t.words {
		if i < len(s.words) {
			s.words[i] &= tWord
			sLen++
		}
	}
	s.words = s.words[:sLen]
}

// DifferenceWith 差集：元素出现在A集合，未出现在B集合
func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, tWord := range t.words {
		if i < len(s.words) {
			s.words[i] &= ^tWord
		} else {
			s.words = append(s.words, tWord)
		}
	}
}

// SymmetricDifference 并差集：元素出现在A但没有出现在B，或者出现在B没有出现在A
func (s *IntSet) SymmetricDifference(t *IntSet) {
	temp := t.Copy()
	temp.DifferenceWith(s)
	s.DifferenceWith(t)
	s.UnionWith(temp)
}

// Elems 返回元素列表
func (s *IntSet) Elems() []int {
	res := make([]int, 0, s.Len())
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < size; j++ {
			if word&(1<<uint(j)) != 0 {
				res = append(res, size*i+j)
			}
		}
	}
	return res
}

// Len return the number of elements
func (s *IntSet) Len() int {
	var res int
	for _, word := range s.words {
		res += ch2.PopCount(uint64(word))
	}
	return res
}

// Remove x from the set
func (s *IntSet) Remove(x int) {
	word, bit := x/size, uint(x%size)
	for word >= len(s.words) {
		return
	}
	s.words[word] &= ^(1 << bit)
}

// Clear remove all elements from the set
func (s *IntSet) Clear() {
	s.words = make([]uint, len(s.words))
}

// Copy return a copy of the set
func (s *IntSet) Copy() *IntSet {
	res := &IntSet{
		words: make([]uint, len(s.words)),
	}
	copy(res.words, s.words)
	return res
}

// AddAll 批量增加
func (s *IntSet) AddAll(elems ...int) {
	for _, elem := range elems {
		s.Add(elem)
	}
}

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < size; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", size*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}
