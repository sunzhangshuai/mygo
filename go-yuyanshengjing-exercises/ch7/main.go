package ch7

// Exercises 练习
type Exercises struct {
}

// Task1 使用来自ByteCounter的思路，实现一个针对单词和行数的计数器。你会发现bufio.ScanWords非常的有用。
func (e *Exercises) Task1() {
	type ByteCounter int

	func (c *ByteCounter) Write(p []byte) (int, error) {
		*c += ByteCounter(len(p)) // convert int to ByteCounter
		return len(p), nil
	}
}

// Task2 写一个带有如下函数签名的函数CountingWriter，传入一个io.Writer接口类型，返回一个把原来的Writer封装在里面的新的Writer类型和一个表示新的写入字节数的int64类型指针。
func (e *Exercises) Task2() {

}

// Task3 为在gopl.io/ch4/treesort（§4.4）中的*tree类型实现一个String方法去展示tree类型的值序列。
func (e *Exercises) Task3() {

}
