package main

type DayMeeting struct {
	Day      int
	Meetings int
}

func bestVacation(daysWithMeetings []DayMeeting, periodLength, vacationLength int) [2]int {
	// 0-based массив встреч
	meetingsPerDay := make([]int, periodLength)
	for _, m := range daysWithMeetings {
		meetingsPerDay[m.Day-1] = m.Meetings
	}

	// сумма встреч в первом окне
	windowSum := 0
	for d := 0; d < vacationLength; d++ {
		windowSum += meetingsPerDay[d]
	}

	minMissed := windowSum
	bestStart := 0 // 0-based

	// скользящее окно по остальным дням
	for start := 1; start <= periodLength-vacationLength; start++ {
		windowSum = windowSum - meetingsPerDay[start-1] + meetingsPerDay[start+vacationLength-1]
		if windowSum < minMissed {
			minMissed = windowSum
			bestStart = start
		}
	}

	// возвращаем 1-based день начала отпуска
	return [2]int{bestStart + 1, minMissed}
}

// func main() {
// 	daysWithMeetings := []DayMeeting{
// 		{Day: 3, Meetings: 1},
// 		{Day: 4, Meetings: 3},
// 		{Day: 14, Meetings: 3},
// 		{Day: 21, Meetings: 3},
// 		{Day: 28, Meetings: 1},
// 	}
// 	periodLength := 30
// 	vacationLength := 7

// 	res := bestVacation(daysWithMeetings, periodLength, vacationLength)
// 	fmt.Println(res) // пример вывода: [5 3]
// }
