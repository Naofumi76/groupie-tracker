package utils

import "strconv"

func GenerateList(n int) []string {
	var list []string
	for i := 1; i <= n; i++ {
		list = append(list, strconv.Itoa(i))
	}
	return list
}
