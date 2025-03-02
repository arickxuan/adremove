package adcore

import (
	"github.com/lqqyt2423/go-mitmproxy/proxy"
	log "github.com/sirupsen/logrus"
)

func NewAdService(filePath string) {
	var config *Config

	if filePath != "" {
		c, err := ParseConfigFile(filePath, nil)
		if err != nil {
			log.Panic(err)
			return
		}
		c.HandlePlugin()
		//log.Println("所有配置：", c)
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

	log.Panic(p.Start())
}
