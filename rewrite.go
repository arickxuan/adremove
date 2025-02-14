package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"regexp"

	"github.com/itchyny/gojq"
	"github.com/tidwall/sjson"
)

func checkUrl(url string, rws []Rewrite) (string, []string) {

	for _, v := range rws {
		regex := regexp.MustCompile(v.Pattern)
		if regex.MatchString(url) {
			return v.Action, v.Params
		}
	}

	return "", []string{}

}

func handleJson(jsonData string, action string, params []string) string {

	if action == "response-body-json-add" {
		jsonData = checkJsonAdd(jsonData, params)
	} else if action == "response-body-json-del" {
		jsonData = checkJsonDel(jsonData, params)
	} else if action == "response-body-json-replace" {
		jsonData = checkJsonReplace(jsonData, params)
	} else if action == "response-body-json-jp" {
		jsonData = checkJsonJq(jsonData, params)
	}else if action == "request-body-json-add" {
		jsonData = checkJsonDel(jsonData, params)
	} else if action == "request-body-json-del" {
		jsonData = checkJsonDel(jsonData, params)
	} else if action == "request-body-json-replace" {
		jsonData = checkJsonReplace(jsonData, params)
	} else if action == "request-body-json-jp" {
		jsonData = checkJsonJq(jsonData, params)
	}

	return jsonData
}

func checkJsonAdd(jsonData string, Params []string) string {
	pp := len(Params) / 2
	var err error
	for i := 0; i < pp; i++ {
		jsonData, err = sjson.Set(jsonData, Params[i], Params[i+1])
		if err != nil {
			return ""
		}
	}

	return jsonData
}

func checkJsonDel(jsonData string, Params []string) string {
	// 删除 age 字段
	var err error
	for _, v := range Params {
		jsonData, err = sjson.Delete(jsonData, v)
		if err != nil {
			return ""
		}
	}
	return jsonData
}

func checkJsonReplace(jsonData string, Params []string) string {
	pp := len(Params) / 2
	var err error
	for i := 0; i < pp; i++ {
		jsonData, err = sjson.Set(jsonData, Params[i], Params[i+1])
		if err != nil {
			return ""
		}
	}

	return jsonData
}

func checkJsonJq(jsonData string, Params []string) string {
	var data interface{}

	// 使用json.Unmarshal解析JSON字符串
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range Params {
		query, err := gojq.Parse(item)
		if err != nil {
			log.Fatalln(err)
		}
		//input := map[string]any{"foo": []any{1, 2, 3}}
		iter := query.Run(data) // or query.RunWithContext
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}
			if err, ok := v.(error); ok {
				if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
					break
				}
				log.Fatalln(err)
			}
			by, err := json.Marshal(v)
			if err != nil {
				break
			}
			jsonData = string(by)
		}
	}
	return jsonData

}

func checkReplaceRegex(jsonData string, Params []string) string {
	return ""

}

func checkMock(jsonData string, Params []string) string {

	return ""
}

func checkRespJsonAdd(jsonData string, Params []string) string {
	pp := len(Params) / 2
	var err error
	for i := 0; i < pp; i++ {
		jsonData, err = sjson.Set(jsonData, Params[i], Params[i+1])
		if err != nil {
			return ""
		}
	}

	return jsonData
}

func checkRespJsonDel(jsonData string, Params []string) string {
	var err error
	for _, v := range Params {
		jsonData, err = sjson.Delete(jsonData, v)
		if err != nil {
			return ""
		}
	}
	return jsonData
}

func checkRespJsonReplace(jsonData string, Params []string) string {
	pp := len(Params) / 2
	var err error
	for i := 0; i < pp; i++ {
		jsonData, err = sjson.Set(jsonData, Params[i], Params[i+1])
		if err != nil {
			return ""
		}
	}

	return jsonData
}

func checkRespJsonJq(jsonData string, Params []string) string {
	var data interface{}

	// 使用json.Unmarshal解析JSON字符串
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range Params {
		query, err := gojq.Parse(item)
		if err != nil {
			log.Fatalln(err)
		}
		//input := map[string]any{"foo": []any{1, 2, 3}}
		iter := query.Run(data) // or query.RunWithContext
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}
			if err, ok := v.(error); ok {
				if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
					break
				}
				log.Fatalln(err)
			}
			by, err := json.Marshal(v)
			if err != nil {
				break
			}
			jsonData = string(by)
		}
	}
	return jsonData
}

func checkRespReplaceRegex(jsonData string, Params []string) string {
	return ""
}

func checkRespMock(jsonData string, Params []string) string {

	return ""
}
