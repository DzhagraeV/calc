package main

import (
	"sort"
)

type Entry struct {
	UserId int
	Steps  int
}

type Champions struct {
	UserIds []int
	Steps   int
}

func findChampions(statistics [][]Entry) Champions {
	days := len(statistics)
	if days == 0 {
		return Champions{UserIds: []int{}, Steps: 0}
	}

	totalSteps := make(map[int]int)  // userId -> total steps across all days
	daysPresent := make(map[int]int) // userId -> number of days the user appeared

	for _, day := range statistics {
		for _, e := range day {
			totalSteps[e.UserId] += e.Steps
			daysPresent[e.UserId]++
		}
	}

	var winners []int
	maxSteps := 0
	for id, cnt := range daysPresent {
		if cnt != days {
			continue // участник пропускал хотя бы один день
		}
		if totalSteps[id] > maxSteps || len(winners) == 0 {
			maxSteps = totalSteps[id]
			winners = []int{id}
		} else if totalSteps[id] == maxSteps {
			winners = append(winners, id)
		}
	}

	sort.Ints(winners)
	return Champions{UserIds: winners, Steps: maxSteps}
}

// func main() {
// 	stats1 := [][]Entry{
// 		{{UserId: 1, Steps: 1000}, {UserId: 2, Steps: 1500}},
// 		{{UserId: 2, Steps: 1000}},
// 	}
// 	fmt.Println(findChampions(stats1)) // {UserIds:[2] Steps:2500}

// 	stats2 := [][]Entry{
// 		{{UserId: 1, Steps: 2000}, {UserId: 2, Steps: 1500}},
// 		{{UserId: 2, Steps: 4000}, {UserId: 1, Steps: 3500}},
// 	}
// 	fmt.Println(findChampions(stats2)) // {UserIds:[1 2] Steps:5500}
// }
