package main

import (
	"fmt"
	"strings"
)

type Node struct {
	Text     string
	Children []*Node
}

// рекурсивная функция для обхода дерева
func printPaths(node *Node, path []string) {
	// добавляем текущий узел в путь
	path = append(path, node.Text)

	if len(node.Children) == 0 {
		// если лист, выводим путь
		fmt.Println(strings.Join(path, " => "))
		return
	}

	// рекурсивно обходим детей
	for _, child := range node.Children {
		printPaths(child, path)
	}
}

// func main() {
// 	// пример дерева
// 	tree := []*Node{
// 		{
// 			Text: "Вещи",
// 			Children: []*Node{
// 				{
// 					Text: "Одежда",
// 					Children: []*Node{
// 						{Text: "Мужская"},
// 						{Text: "Женская"},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			Text: "Хобби",
// 			Children: []*Node{
// 				{
// 					Text: "Велосипеды",
// 					Children: []*Node{
// 						{Text: "Горные"},
// 					},
// 				},
// 				{Text: "Мангалы"},
// 			},
// 		},
// 		{
// 			Text: "Транспорт",
// 		},
// 	}

// 	// обходим корневые узлы
// 	for _, node := range tree {
// 		printPaths(node, []string{})
// 	}
// }
