package ch6

import (
	"fmt"
)

// Exercises 练习
type Exercises struct {
}

// Task1
// 1. 为bit数组实现下面这些方法。 Len、Remove、Clear、Copy
// 2. 定义一个变参方法(*IntSet).AddAll(...int)，这个方法可以添加一组IntSet，比如s.AddAll(1,2,3)。
// 3. (*IntSet).UnionWith会用|操作符计算两个集合的并集，我们再为IntSet实现另外的几个函数：
//		IntersectWith（交集：元素在A集合B集合均出现），
//		DifferenceWith（差集：元素出现在A集合，未出现在B集合），
//		SymmetricDifference（并差集：元素出现在A但没有出现在B，或者出现在B没有出现在A）。
// 4. 实现一个Elems方法，返回集合中的所有元素，用于做一些range之类的遍历操作。
// 5. 我们这章定义的IntSet里的每个字都是用的uint64类型，但是64位的数值可能在32位的平台上不高效。修改程序，使其使用uint类型，这种类型对于32位平台来说更合适。当然了，这里我们可以不用简单粗暴地除64，可以定义一个常量来决定是用32还是64，这里你可能会用到平台的自动判断的一个智能表达式：32 << (^uint(0) >> 63)
// go run main.go -ch=6 -task=1
func (e *Exercises) Task1() {
	var x, y IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	fmt.Println("x elems：", x.String()) // "{1 9 144}"

	y.Add(9)
	y.Add(42)
	fmt.Println("y elems：", y.String()) // "{9 42}"

	// 交、差、并差
	x1 := x.Copy()
	x1.IntersectWith(&y)
	fmt.Println("x1 elems：", x1.String()) // "{9}"
	x2 := x.Copy()
	x2.DifferenceWith(&y)
	fmt.Println("x2 elems：", x2.String()) // "{1 144}"
	x3 := x.Copy()
	x3.SymmetricDifference(&y)
	fmt.Println("x3 elems：", x3.String()) // "{1 42 144}"

	x.UnionWith(&y)

	z := x.Copy()
	fmt.Println("z elems：", z.String()) // "{1 9 42 144}"

	z.AddAll(7, 8, 9)
	fmt.Println("z elems：", z.String()) // "{1 7 8 9 42 144}"

	fmt.Println("z elems：", z.Elems()) // "{1 7 8 9 42 144}"

	x.Remove(9)
	fmt.Println("x elems：", x.String())                         // "{1 42 144}"
	fmt.Println("x has 9：", x.Has(9), "x has 123：", x.Has(123)) // "false false"
	fmt.Println("x len：", x.Len())                              // "3"
	x.Clear()
	fmt.Println("x elems：", x.String()) // "{}"
}

// Task2 to Task1
func (e *Exercises) Task2() {
	e.Task1()
}

// Task3 to Task1
func (e *Exercises) Task3() {
	e.Task1()
}

// Task4 to Task1
func (e *Exercises) Task4() {
	e.Task1()
}

// Task5 to Task1
func (e *Exercises) Task5() {
	e.Task1()
}
