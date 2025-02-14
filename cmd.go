package main

import (
	"os"

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
	var config *Config

	if args[1] != "" {
		c, err := ParseConfigFile(args[1], nil)
		if err != nil {
			log.Fatal(err)
			return
		}
		c.HandlePlugin()
		log.Println("所有配置：", c)
		config = c

	} else {
		config = &Config{}
		config.Addr = ":9080"
		config.SslInsecure = false
	}

	opts := &proxy.Options{
		Addr:              config.Addr,
		StreamLargeBodies: 1024 * 1024 * 5,
		SslInsecure:       config.SslInsecure,
	}
	if config.CaRootPath != "" {
		opts.CaRootPath = config.CaRootPath
	}
	if config.EnableCustomCa {
		opts.NewCaFunc = NewTrustedCA
	}

	p, err := proxy.NewProxy(opts)
	if err != nil {
		log.Fatal(err)
	}
	add := &CloseConn{
		config: config,
	}
	p.AddAddon(add)
	p.AddAddon(&proxy.LogAddon{})

	log.Fatal(p.Start())
}
