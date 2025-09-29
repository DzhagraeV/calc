package main

import (
	"sort"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// lowerBound возвращает первый индекс i в a, такой что a[i] >= x.
// Если такого нет — возвращает len(a).
func lowerBound(a []int, x int) int {
	if len(a) == 0 {
		return 0
	}

	l := 0
	r := len(a) - 1

	for l <= r {
		m := l + (r-l)/2
		if a[m] < x {
			l = m + 1
		} else {
			r = m - 1
		}
	}

	return l
}

func FindGood(goods []int, needs []int) int {
	if len(goods) == 0 {
		sum := 0
		for _, need := range needs {
			sum += abs(need)
		}
		return sum
	}

	sort.Ints(goods)
	sum := 0

	for _, need := range needs {
		idx := lowerBound(goods, need)

		if idx == 0 {
			sum += abs(goods[0] - need)
		} else if idx == len(goods) {
			sum += abs(goods[len(goods)-1] - need)
		} else {
			leftDiff := abs(goods[idx-1] - need)
			rightDiff := abs(goods[idx] - need)
			if leftDiff <= rightDiff {
				sum += leftDiff
			} else {
				sum += rightDiff
			}
		}
	}

	return sum
}

// func main() {
// 	goods := []int{8, 3, 5}
// 	needs := []int{5, 14, 12, 44, 55}
// 	fmt.Println(FindGood(goods, needs)) // 93
// }
