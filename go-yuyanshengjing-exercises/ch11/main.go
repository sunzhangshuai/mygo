package ch11

import (
	"bufio"
	"io"
	"math/rand"
	"unicode"
)

// charCount 字符统计
func charCount(char io.Reader) map[rune]int {
	counts := make(map[rune]int) // counts of Unicode characters
	invalid := 0

	in := bufio.NewReader(char)
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
	}
	return counts
}

// IsPalindrome reports whether s reads the same forward and backward.
// (Our first attempt.)
func IsPalindrome(s string) bool {
	r := make([]rune, 0, len(s))
	for _, c := range s {
		r = append(r, c)
	}
	for i := 0; i < len(r)/2; i++ {
		if r[i] != r[len(r)-1-i] {
			return false
		}
	}
	return true
}

// randomPalindrome returns a palindrome whose length and contents
// are derived from the pseudo-random number generator rng.
func randomPalindrome(rng *rand.Rand) string {
	n := rng.Intn(25) // random length up to 24
	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		r := rune(rng.Intn(0x1000)) // random rune up to '\u0999'
		runes[i] = r
		runes[n-1-i] = r
	}
	return string(runes)
}

// notRandomPalindrome 非回文串
func notRandomPalindrome(rng *rand.Rand) string {
	n := rng.Intn(25) // random length up to 24
	n += 2
	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		r := rune(rng.Intn(0x1000)) // random rune up to '\u0999'
		runes[i] = r
		for v := rune(rng.Intn(0x1000)); v != r; v = rune(rng.Intn(0x1000)) {
			runes[n-1-i] = v
		}
	}
	return string(runes)
}

// pc[i] is the population count of i.
var pc = func() (pc [256]byte) {
	for i := range pc {
		// 右移一位的数字已经统计过了，看看最低位是不是1。
		pc[i] = pc[i>>1] + byte(i&1)
	}
	return
}()

// PopCountTable returns the population count (number of set bits) of x.
func PopCountTable(x uint64) int {
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

// PopCountBit 位移法计算
func PopCountBit(x uint64) int {
	var result int
	for ; x != 0; x = x >> 1 {
		if x&1 == 1 {
			result++
		}
	}
	return result
}

// PopCountBitAnd 位与法
func PopCountBitAnd(x uint64) int {
	var result int
	for x > 0 {
		x = x & (x - 1)
		result++
	}
	return result
}
