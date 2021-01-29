package common

import "strings"

func Min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}


func GetPath(currentGroup, relativePath string) string {
	g := "/" + currentGroup
	if g == "/" {
		g = ""
	}
	g = g + relativePath
	g = strings.Replace(g, "//", "/", -1)
	return g
}
