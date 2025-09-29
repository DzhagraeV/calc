package main

import "fmt"

func visit(node string, deps map[string][]string, seen map[string]bool) {
	if seen[node] {
		return
	}
	seen[node] = true
	fmt.Println(node) // обработка узла
	for _, dep := range deps[node] {
		visit(dep, deps, seen)
	}
}

func dep() {
	deps := map[string][]string{
		"tensorflow": {"nvcc", "gpu", "linux"},
		"nvcc":       {"linux"},
		"linux":      {"core"},
		"mylib":      {"tensorflow"},
		"mylib2":     {"requests"},
	}

	seen := make(map[string]bool)
	for node := range deps {
		visit(node, deps, seen)
	}
}
