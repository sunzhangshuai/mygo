package ch2

import "testing"

var testNum uint64 = 100000000

func BenchmarkExercises_PopCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCount(testNum)
	}
}

func BenchmarkExercises_Task3(b *testing.B) {
	e := Exercises{}
	for i := 0; i < b.N; i++ {
		e.Task3(testNum)
	}
}

func BenchmarkExercises_Task4(b *testing.B) {
	e := Exercises{}
	for i := 0; i < b.N; i++ {
		e.Task4(testNum)
	}
}

func BenchmarkExercises_Task5(b *testing.B) {
	e := Exercises{}
	for i := 0; i < b.N; i++ {
		e.Task5(testNum)
	}
}
