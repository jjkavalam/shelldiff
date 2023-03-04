package shelldiff

func longestCommonSubsequence[T any](s1, s2 []T, compareFn func(a T, b T) bool) []T {
	m, n := len(s1), len(s2)
	if m == 0 || n == 0 {
		return make([]T, 0)
	}
	memo := make([][]int, m+1)
	for i := range memo {
		memo[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if compareFn(s1[i-1], s2[j-1]) {
				memo[i][j] = memo[i-1][j-1] + 1
			} else {
				memo[i][j] = max(memo[i-1][j], memo[i][j-1])
			}
		}
	}
	lcsLen := memo[m][n]
	if lcsLen == 0 {
		return make([]T, 0)
	}
	lcs := make([]T, lcsLen)
	i, j := m, n
	for lcsLen > 0 {
		if compareFn(s1[i-1], s2[j-1]) {
			lcs[lcsLen-1] = s1[i-1]
			i--
			j--
			lcsLen--
		} else if memo[i-1][j] > memo[i][j-1] {
			i--
		} else {
			j--
		}
	}
	return lcs
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
