package main

import (
	"fmt"
)

type Ticket struct {
	From string
	To   string
}

func ticket() {
	tickets := []Ticket{
		{From: "London", To: "Moscow"},
		{From: "NY", To: "London"},
		{From: "Moscow", To: "SPb"},
	}

	// map from -> to
	fromTo := make(map[string]string, len(tickets))
	// set of "to" городов
	toSet := make(map[string]struct{}, len(tickets))

	for _, t := range tickets {
		fromTo[t.From] = t.To
		toSet[t.To] = struct{}{}
	}

	// найти начало — город, который встречается как From, но не встречается в toSet
	start := ""
	for _, t := range tickets {
		if _, ok := toSet[t.From]; !ok {
			start = t.From
			break
		}
	}

	// собираем маршрут без visited
	route := make([]Ticket, 0, len(tickets))
	cur := start
	for {
		to, ok := fromTo[cur]
		if !ok {
			break // дошли до конца маршрута
		}
		route = append(route, Ticket{From: cur, To: to})
		cur = to
	}

	// вывод
	fmt.Println("route := []Ticket{")
	for _, r := range route {
		fmt.Printf("    {From: %q, To: %q},\n", r.From, r.To)
	}
	fmt.Println("}")
}
