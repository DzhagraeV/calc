package main

import (
	"strconv"
	"strings"
)

func compareVersions(v1, v2 string) int {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	minLen := len(parts1)
	if len(parts2) < minLen {
		minLen = len(parts2)
	}

	// 1. сравниваем части до минимальной длины
	for i := 0; i < minLen; i++ {
		num1, _ := strconv.Atoi(parts1[i])
		num2, _ := strconv.Atoi(parts2[i])
		if num1 < num2 {
			return 1 // v2 > v1
		} else if num1 > num2 {
			return -1 // v1 > v2
		}
	}

	// 2. проверяем остаток более длинного массива
	if len(parts1) > minLen {
		for i := minLen; i < len(parts1); i++ {
			num, _ := strconv.Atoi(parts1[i])
			if num > 0 {
				return -1 // v1 > v2
			}
		}
	} else if len(parts2) > minLen {
		for i := minLen; i < len(parts2); i++ {
			num, _ := strconv.Atoi(parts2[i])
			if num > 0 {
				return 1 // v2 > v1
			}
		}
	}

	return 0 // равны
}

// func main() {
// 	fmt.Println(compareVersions("v11.22.44", "v11.22.45"))   // 1
// 	fmt.Println(compareVersions("v11.22.44", "v11.22.44"))   // 0
// 	fmt.Println(compareVersions("v11.22.44", "v11.22.44.0")) // 0
// 	fmt.Println(compareVersions("v1.12.3", "v1.3.4"))        // -1
// }
