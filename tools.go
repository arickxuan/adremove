package main

import (
	"regexp"
	"strings"
)

func isInList(str string, list []string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

func extractContent(s string) string {
	start := strings.Index(s, "(")   // 找到第一个 (
	end := strings.LastIndex(s, ")") // 找到最后一个 )

	if start == -1 || end == -1 || start >= end {
		return "" // 如果找不到或者位置不正确，则返回空字符串
	}
	return s[start+1 : end] // 提取括号中的内容
}

func trimBlank(str string) string {
	regex := regexp.MustCompile(`\s+`)
	noSpaceStr := regex.ReplaceAllString(str, "")
	return noSpaceStr
}
