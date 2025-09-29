package main

func addBigNumbers(arr1, arr2 []int) []int {
	i, j := len(arr1)-1, len(arr2)-1
	overflow := false
	result := []int{}

	for i >= 0 || j >= 0 {
		sum := 0
		if i >= 0 {
			sum += arr1[i]
			i--
		}
		if j >= 0 {
			sum += arr2[j]
			j--
		}
		if overflow {
			sum += 1
		}

		if sum >= 10 {
			sum -= 10
			overflow = true
		} else {
			overflow = false
		}

		result = append([]int{sum}, result...)
	}

	// добавляем overflow после основного цикла, если он остался
	if overflow {
		result = append([]int{1}, result...)
	}

	return result
}

// func main() {
// 	arr1 := []int{1, 2, 3}
// 	arr2 := []int{4, 5, 6}
// 	fmt.Println(addBigNumbers(arr1, arr2)) // [5 7 9]

// 	arr3 := []int{5, 4, 4}
// 	arr4 := []int{4, 5, 6}
// 	fmt.Println(addBigNumbers(arr3, arr4)) // [1 0 0 0]
// }
