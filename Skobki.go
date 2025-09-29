package main

func generateParenthesis(n int) []string {
	if n == 0 {
		return []string{""}
	}

	res := make([]string, 0)
	buf := make([]byte, 2*n) // буфер для текущей последовательности

	var backtrack func(pos, open, close int)
	backtrack = func(pos, open, close int) {
		if pos == 2*n {
			res = append(res, string(buf))
			return
		}
		// можно добавить '(' если ещё не использовали все открывающие
		if open < n {
			buf[pos] = '('
			backtrack(pos+1, open+1, close)
		}
		// можно добавить ')' только если есть незакрытые '('
		if close < open {
			buf[pos] = ')'
			backtrack(pos+1, open, close+1)
		}
	}

	backtrack(0, 0, 0)
	return res
}

// func main() {
// 	fmt.Println(generateParenthesis(3))
// }
