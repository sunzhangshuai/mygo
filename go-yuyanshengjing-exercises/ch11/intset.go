package ch11

import (
	"bytes"
	"sort"
	"strconv"
)

// IntSet 整数集合
type IntSet struct {
	words map[int]bool
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	_, ok := s.words[x]
	return ok
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	if s.words == nil {
		s.words = make(map[int]bool)
	}
	s.words[x] = true
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for k := range t.words {
		s.words[k] = true
	}
}

// IntersectWith 交集：元素在A集合B集合均出现
func (s *IntSet) IntersectWith(t *IntSet) {
	for k := range s.words {
		if _, ok := t.words[k]; !ok {
			delete(s.words, k)
		}
	}
}

// DifferenceWith 差集：元素出现在A集合，未出现在B集合
func (s *IntSet) DifferenceWith(t *IntSet) {
	for k := range s.words {
		if _, ok := t.words[k]; ok {
			delete(s.words, k)
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
	for k := range s.words {
		res = append(res, k)
	}
	return res
}

// Len return the number of elements
func (s *IntSet) Len() int {
	return len(s.words)
}

// Remove x from the set
func (s *IntSet) Remove(x int) {
	delete(s.words, x)
}

// Clear remove all elements from the set
func (s *IntSet) Clear() {
	s.words = make(map[int]bool)
}

// Copy return a copy of the set
func (s *IntSet) Copy() *IntSet {
	res := &IntSet{
		words: make(map[int]bool, len(s.words)),
	}
	for k := range s.words {
		res.words[k] = true
	}
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

	intList := sort.IntSlice(s.Elems())
	sort.Sort(intList)
	for i := 0; i < s.Len(); i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(strconv.Itoa(intList[i]))
	}
	buf.WriteByte('}')
	return buf.String()
}
