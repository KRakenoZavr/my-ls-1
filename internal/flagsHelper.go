package internal

import (
	"sort"
	"strings"
)

func RemoveDotFiles(s []fileInfo) []fileInfo {
	files := []fileInfo{}
	for _, l := range s {
		if !strings.HasPrefix(l.name, ".") {
			files = append(files, l)
		}
	}

	return files
}

func SortByTime(s []fileInfo) {
	sort.SliceStable(s, func(i, j int) bool {
		return s[i].fullDate.After(s[j].fullDate)
	})
}

func ReverseArray(s []fileInfo) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
