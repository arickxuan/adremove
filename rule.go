package main

import (
	"log"
	"net/url"
	"regexp"
	"strings"
)

type RuleItemData struct {
	Type string
	Data string
}

func CheckDOMAIN(url string, target string) bool {
	//return
	return url == target
}

// 域名关键词匹配
func CheckDOMAINKEYWORD(url string, target string) bool {
	//使用 strings.Contains 函数检查 url 中是否包含 target
	return strings.Contains(url, target)
}

func CheckDOMAINSUFFIX(url string, target string) bool {
	// 使用 strings.HasSuffix 函数检查 url 是否以 target 结尾
	return strings.HasSuffix(url, target)
}
func CheckURLREGEX(url string, target string) bool {
	regex := regexp.MustCompile(target)
	return regex.MatchString(url)
}

func CheckAdd(f *url.URL, target string) bool {
	regex := regexp.MustCompile(`\(([^,]+),([^\)]+)\)`)
	matches := regex.FindAllStringSubmatch(extractContent(target), -1)

	keys := make([]RuleItemData, len(matches))

	for i, match := range matches {
		keys[i].Type = trimBlank(match[1])
		keys[i].Data = trimBlank(match[2])
	}
	log.Println(keys)
	for _, v := range keys {
		if v.Type == "DOMAIN-KEYWORD" {
			if !CheckDOMAINKEYWORD(f.Host, v.Data) {
				return false
			}
		} else if v.Type == "DOMAIN-SUFFIX" {
			if !CheckDOMAINSUFFIX(f.Host, v.Data) {
				return false
			}
		} else if v.Type == "DOMAIN" {
			if !CheckDOMAIN(f.Host, v.Data) {
				return false
			}
		} else if v.Type == "URL-REGEX" {
			if !CheckURLREGEX(f.String(), v.Data) {
				return false
			}
		}

	}

	return true
}

func CheckRule(f *url.URL, rules []Rule) (bool, string) {

	for _, v := range rules {
		if v.Type == "DOMAIN" {
			if CheckDOMAIN(f.Host, v.Params) {
				return true, v.Action
			}
		} else if v.Type == "DOMAIN-KEYWORD" {
			if CheckDOMAINKEYWORD(f.Host, v.Params) {
				return true, v.Action
			}
		} else if v.Type == "DOMAIN-SUFFIX" {
			if CheckDOMAINSUFFIX(f.Host, v.Params) {
				return true, v.Action
			}
		} else if v.Type == "URL-REGEX" {
			if CheckURLREGEX(f.String(), v.Params) {
				return true, v.Action
			}
		} else if v.Type == "ADD" {
			if CheckAdd(f, v.Params) {
				return true, v.Action
			}

		}

	}
	return false, ""

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
