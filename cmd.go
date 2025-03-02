package main

import (
	"os"

	"adremove/adcore"

	"github.com/lqqyt2423/go-mitmproxy/proxy"
	log "github.com/sirupsen/logrus"
)

func back(f *proxy.Flow) {

	//修改请求头
	log.Printf("Host: %v, Method: %v, Scheme: %v", f.Request.URL.Host, f.Request.Method, f.Request.URL.Scheme)
	f.Request.URL.Host = "www.baidu.com"
	f.Request.URL.Scheme = "http"
	log.Printf("After: %v", f.Request.URL)

	f.Done()
	resp := &proxy.Response{
		StatusCode: 200,
		Header:     nil,
		Body:       []byte("changed response"),
		BodyReader: nil,
	}
	f.Response = resp
}

func main() {
	// 获取所有命令行参数
	args := os.Args
	// 打印所有参数
	log.Println("所有参数：", args)

	adcore.NewAdService(args[1])
}
