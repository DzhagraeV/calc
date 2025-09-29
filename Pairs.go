package main

func findPairs(nums []int, target int) [][]int {
	pairs := [][]int{}
	seen := make(map[int]bool)

	for _, num := range nums {
		complement := target - num
		if seen[complement] {
			pairs = append(pairs, []int{complement, num})
		}
		seen[num] = true
	}

	return pairs
}

// func main() {
// 	nums := []int{2, 4, 5, 3}
// 	target := 7

// 	result := findPairs(nums, target)
// 	fmt.Println(result) // [[2 5] [4 3]]
// }
