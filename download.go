package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func getFileName(url string) string {
	// 使用正则表达式匹配并替换非字母和数字的字符

	// 去掉开头的 http:// 或 https://
	processedStr := strings.TrimPrefix(url, "http://")
	processedStr = strings.TrimPrefix(processedStr, "https://")

	reg := regexp.MustCompile("[^a-zA-Z0-9]")
	return reg.ReplaceAllString(processedStr, "") + ".plugin"

}

func downLoad(url string, fileName string) {
	// 设置自定义的 User-Agent
	userAgent := "Mozilla/5.0 (iPhone; CPU iPhone OS 16_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Mobile/15E148 Safari/604.1"

	// 创建一个 HTTP 客户端
	client := &http.Client{}

	// 创建一个 GET 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// 设置自定义的 User-Agent
	req.Header.Set("User-Agent", userAgent)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 创建一个文件来保存下载的内容
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// 将响应的内容写入文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("File downloaded successfully.")
}
