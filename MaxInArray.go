package main

// возвращает индекс минимального элемента в срезе
func minIndex(arr []int) int {
	minIdx := 0
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[minIdx] {
			minIdx = i
		}
	}
	return minIdx
}

func topKBuffer(nums []int, k int) []int {
	// 1. первые K элементов в буфер
	buffer := make([]int, k)
	copy(buffer, nums[:k])

	// 2. обрабатываем оставшиеся элементы
	for _, num := range nums[k:] {
		minIdx := minIndex(buffer)
		if num > buffer[minIdx] {
			buffer[minIdx] = num
		}
	}

	return buffer
}

// func main() {
// 	nums := []int{100, 50, 0, 150, 100, 0, -30, 70}
// 	k := 3

// 	result := topKBuffer(nums, k)
// 	fmt.Println(result) // возможный вывод: [100 150 100]
// }
