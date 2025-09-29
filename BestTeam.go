package main

// Struct для команды
type Team struct {
	Backend  int
	Frontend int
	QA       int
	Design   int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// bestTeam собирает команду с минимальной разницей уровней
func bestTeam(backend, frontend, qa, design []int) Team {
	i, j, k, l := 0, 0, 0, 0
	bestRange := int(^uint(0) >> 1) // максимальное int
	var result Team

	for i < len(backend) && j < len(frontend) && k < len(qa) && l < len(design) {
		// берем текущие значения
		b := backend[i]
		f := frontend[j]
		q := qa[k]
		d := design[l]

		currentMin := min(min(b, f), min(q, d))
		currentMax := max(max(b, f), max(q, d))
		currentRange := currentMax - currentMin

		if currentRange < bestRange {
			bestRange = currentRange
			result = Team{Backend: b, Frontend: f, QA: q, Design: d}
		}

		// сдвигаем указатель на минимальный элемент
		if currentMin == b {
			i++
		} else if currentMin == f {
			j++
		} else if currentMin == q {
			k++
		} else {
			l++
		}
	}

	return result
}

// func main() {
// 	// пример 1
// 	backend := []int{1, 2, 2, 3}
// 	frontend := []int{1, 3}
// 	qa := []int{3, 4, 4}
// 	design := []int{2, 3}

// 	team := bestTeam(backend, frontend, qa, design)
// 	fmt.Println(team) // {3 3 3 3}

// 	// пример 2
// 	backend2 := []int{5}
// 	frontend2 := []int{3, 6, 7, 10}
// 	qa2 := []int{3, 9, 11, 18}
// 	design2 := []int{20}

// 	team2 := bestTeam(backend2, frontend2, qa2, design2)
// 	fmt.Println(team2) // {5 6 9 20}
// }
