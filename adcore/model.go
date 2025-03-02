package adcore

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

type General struct {
	Addr           string `json:"addr"`
	CaRootPath     string `json:"ca_root_path"`
	EnableCustomCa bool   `json:"enable_custom_ca"`
	SslInsecure    bool   `json:"ssl_insecure"`
}

// Plugin represents a plugin entry in the configuration
type Plugin struct {
	URL     string `json:"url"`
	Enabled bool   `json:"enabled"`
}

// Rule represents a filtering rule
type Rule struct {
	Params string `json:"params"`
	Action string `json:"action"`
	Type   string `json:"type"`
}

type RuleItem struct {
}

// Rewrite represents a rewrite rule
type Rewrite struct {
	Pattern string   `json:"pattern"`
	Action  string   `json:"action"`
	Params  []string `json:"params"`
}

// MitM represents the MitM hostname list
type MitM struct {
	Hostnames []string `json:"hostnames"`
}

// Config represents the entire parsed configuration
type Config struct {
	General
	Plugins  []Plugin  `json:"plugins"`
	Rules    []Rule    `json:"rules"`
	Rewrites []Rewrite `json:"rewrites"`
	MitM     MitM      `json:"mitm"`
}

func (c *Config) HandlePlugin() {
	for _, p := range c.Plugins {
		file := getFileName(p.URL)
		log.Println(file)
		downLoad(p.URL, file)
		ParseConfigFile(file, c)
	}

}

func (c *Config) HandleGeneral() {

}

func (c *Config) HandleRules() {
	for _, v := range c.Rules {
		if v.Type == "rule" {
			// 处理规则
			// 示例：根据规则类型执行相应的操作
			switch v.Action {
			case "reject":
				// 执行拒绝操作
			case "response-body-json-del":
				// 执行删除操作
			}

		}
	}

}

// NewConfig creates and initializes a Config instance
func NewConfig() *Config {
	return &Config{
		Plugins: []Plugin{
			{URL: "https://kelee.one/Tool/Loon/Plugin/BoxJs.plugin", Enabled: true},
			{URL: "https://kelee.one/Tool/Loon/Plugin/TikTok_redirect.plugin", Enabled: false},
		},
		Rewrites: []Rewrite{
			{Pattern: `^https:\/\/www\.123pan\.com\/api\/config\/get`, Action: "reject-dict"},
			{Pattern: `^https:\/\/www\.123pan\.com\/home`, Action: "reject"},
			{Pattern: `^https:\/\/xiaoshuo\.wtzw\.com\/api\/v\d\/user\/my-center\?`, Action: "response-body-json-del data.func_area[1] data.func_area[2]"},
		},
		MitM: MitM{
			Hostnames: []string{"www.123pan.com", "video-dsp.pddpic.com", "t-dsp.pinduoduo.com", "images.pinduoduo.com"},
		},
	}
}

// ParseConfigFile parses a given configuration file into a Config struct
func ParseConfigFile(filename string, config *Config) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if config == nil {
		config = &Config{}
	}
	var currentSection string
	pluginRegex := regexp.MustCompile(`^(https?://\S+)\s*,\s*enabled\s*=\s*(true|false)\s*$`)
	ruleRegex := regexp.MustCompile(`^([^,]+)\s*,\s*((?:[^,]+,)*[^,]+)\s*,\s*([^,]+)$`)
	//rewriteRegex := regexp.MustCompile(`^(\S+)(?:\s+(\S+))+$`)
	mitmRegex := regexp.MustCompile(`^hostname\s*=\s*(.+)$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.ToLower(strings.Trim(line, "[]"))
			continue
		}

		switch currentSection {
		case "general":
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				switch key {
				case "addr":
					config.Addr = value
				case "ca_root_path":
					config.CaRootPath = value
				case "enable_custom_ca":
					config.EnableCustomCa = value == "true"
				case "ssl_insecure":
					config.SslInsecure = value == "true"
				}
			}

		case "plugin":
			if matches := pluginRegex.FindStringSubmatch(line); matches != nil {
				log.Println("matches", matches)
				enabled := trimBlank(matches[2]) == "true"
				log.Println("matches2", trimBlank(matches[2]))
				log.Println("matches1", trimBlank(matches[1]))
				if enabled {
					config.Plugins = append(config.Plugins, Plugin{URL: matches[1], Enabled: enabled})
				}
			}
		case "rule":
			if matches := ruleRegex.FindStringSubmatch(line); matches != nil {
				params := trimBlank(matches[2])
				ty := trimBlank(strings.ToUpper(matches[1]))
				config.Rules = append(config.Rules, Rule{Type: ty, Params: params, Action: trimBlank(matches[3])})
			}
		case "rewrite":
			// if matches := rewriteRegex.FindStringSubmatch(line); matches != nil {
			// 	log.Println("matches[0]", matches[0])
			// 	config.Rewrites = append(config.Rewrites, Rewrite{Pattern: matches[1], Action: matches[2], Params: matches[3:]})
			// }
			parts := strings.Split(line, " ")
			if len(parts) >= 2 {
				config.Rewrites = append(config.Rewrites, Rewrite{Pattern: parts[0], Action: strings.ToLower(parts[1]), Params: parts[2:]})
			}

		case "mitm":
			if matches := mitmRegex.FindStringSubmatch(line); matches != nil {
				hostnames := strings.Split(matches[1], ", ")
				config.MitM.Hostnames = append(config.MitM.Hostnames, hostnames...)
			}
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

//{Pattern: regexp.MustCompile(`^http:\/\/pc\.suishenyun\.net\/peacock\/api\/adspool\?`), Action: "response-body-json-jq '.data |= map((.layout |= map(select(.key_name == \"实用工具\" or .key_name == \"会员下icon\" or .key_name == \"天气页面_右上角图标\" or .key_name == \"天气页面_24小时\" or .key_name == \"天气页面_15日\" or .key_name == \"每日插屏\" or .key_name == \"下拉屏保\" or .key_name == \"星座Tab（黄历页）"))) | select(.layout | length > 0))'"},
