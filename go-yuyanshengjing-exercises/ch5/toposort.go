package ch5

import (
	"fmt"
)

// preReqs记录了每个课程的前置课程
var preReqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

const (
	notVisit = iota
	visiting
	visited
)

// topologySort 拓扑排序
func topologySort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]int)
	var visitAll func(items map[string]bool)
	var hasRing bool
	visitAll = func(items map[string]bool) {
		if hasRing {
			return
		}
		for item := range items {
			if seen[item] == notVisit {
				seen[item] = visiting
				keys := make(map[string]bool)
				for _, key := range m[item] {
					keys[key] = true
				}
				visitAll(keys)
				seen[item] = visited
				order = append(order, item)
			} else if seen[item] == visiting {
				hasRing = true
				return
			}
		}
	}
	keys := make(map[string]bool)
	for key := range m {
		keys[key] = true
	}
	visitAll(keys)
	if hasRing {
		fmt.Println("拓扑排序有环")
		return nil
	}
	return order
}
