package adcore

import (
	"net/url"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/lqqyt2423/go-mitmproxy/proxy"
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

func AdResponse(action string) *proxy.Response {
	resp := &proxy.Response{}

	if action == "reject" { // 直接断开连接
		resp.StatusCode = 502
		resp.Header.Set("Content-Type", "text/plain")
	} else if action == "reject-200" { // 返回一个code为200，body内容为空的response
		resp.StatusCode = 200
		resp.Header.Set("Content-Type", "text/plain")
		resp.Body = []byte("")
	} else if action == "reject_img" { //返回一个code为200，body内容一像素图片的的response
		resp.StatusCode = 200
		resp.Header.Set("Content-Type", "image/png")
		resp.Body = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x02\x00\x00\x00\x90wS\xde\x00\x00\x00\x0cIDAT\x08\xd7c\xf8\xff\xff?\x00\x05\xfe\x02\xfe\xdc\xcc\x11\x00\x00\x00\x00IEND\xaeB`\x82")
	} else if action == "reject_dict" { //返回一个code为200，body内容为"{}"的空json对象字符串
		resp.StatusCode = 200
		resp.Header.Set("Content-Type", "application/json")
		resp.Body = []byte("{}")
	} else if action == "reject_array" { //返回一个code为200，body内容为"[]"的空json数组字符串
		resp.StatusCode = 200
		resp.Header.Set("Content-Type", "application/json")
		resp.Body = []byte("[]")
	} else if action == "reject_str" { //返回一个code为200，body内容为"reject_str"的字符串
		resp.StatusCode = 200
		resp.Body = []byte("reject_str")
	}

	return resp
}
