package main

func moveZerosRight(nums []int) {
	pos := 0 // позиция для следующего ненулевого элемента

	// сначала сдвигаем все ненулевые числа влево
	for _, num := range nums {
		if num != 0 {
			nums[pos] = num
			pos++
		}
	}

	// заполняем оставшиеся позиции нулями
	for pos < len(nums) {
		nums[pos] = 0
		pos++
	}
}

// func main() {
// 	arr := []int{7, 3, 0, 0, 0, 2, 4, 0, 5, 19}
// 	moveZerosRight(arr)
// 	fmt.Println(arr) // [7 3 2 4 5 19 0 0 0 0]
// }
