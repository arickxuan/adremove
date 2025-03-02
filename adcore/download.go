package adcore

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
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
	userAgent := "Loon/821 CFNetwork/1399 Darwin/22.1.0"

	// 创建一个 HTTP 客户端
	client := &http.Client{}

	// 创建一个 GET 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}

	// 设置自定义的 User-Agent
	req.Header.Set("User-Agent", userAgent)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 创建一个文件来保存下载的内容
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// 将响应的内容写入文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println("Error writing to file:", err)
		return
	}

	log.Println("File downloaded successfully.")
}
