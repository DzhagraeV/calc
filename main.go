package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println(lowerBound([]int{3, 5, 8}, 5))

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	rand.Seed(time.Now().Unix())

	r.Int()
}
