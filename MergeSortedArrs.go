package main

func mergeSortedArrays(A, B []int) []int {
	M, N := len(A), len(B)
	result := make([]int, 0, M+N)

	i, j := 0, 0
	for i < M && j < N {
		if A[i] <= B[j] {
			result = append(result, A[i])
			i++
		} else {
			result = append(result, B[j])
			j++
		}
	}

	// добавляем оставшиеся элементы
	for i < M {
		result = append(result, A[i])
		i++
	}
	for j < N {
		result = append(result, B[j])
		j++
	}

	return result
}

// func main() {
// 	A1 := []int{1, 2, 5}
// 	B1 := []int{1, 2, 3, 4, 6}
// 	fmt.Println(mergeSortedArrays(A1, B1)) // [1 1 2 2 3 4 5 6]

// 	A2 := []int{4, 7, 13}
// 	B2 := []int{3, 5, 8, 9, 11}
// 	fmt.Println(mergeSortedArrays(A2, B2)) // [3 4 5 7 8 9 11 13]
// }
