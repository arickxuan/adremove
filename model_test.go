package main

import (
	"fmt"
	"log"
	"net/url"
	"testing"

	"github.com/itchyny/gojq"
)

func TestP12(t *testing.T) {
	hashCa()
	//base64P12 := "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"
	//base64P12ToP12File(base64P12)
}

func TestConfig(t *testing.T) {
	c, err := ParseConfigFile("rule.conf", nil)
	if err != nil {
		t.Log(err)
	}
	//t.Log(c)
	t.Log("ok")
	for _, v := range c.Plugins {
		t.Log(v.Enabled, v.URL)
	}

	for _, v := range c.Rules {
		t.Log(v.Type, v.Action, v.Params)
	}
}

func TestRule(t *testing.T) {
	t.Log("start")
	str := "https://www.marscode.cn/ide/4633d2l64x7ovd"
	u, err := url.Parse(str)
	if err != nil {
		t.Log(err)
	}
	t.Log("url:", u.String())
	re := CheckAdd(u, " (  (   DOMAIN-KEYWORD , -ad- ), (  DOMAIN-SUFFIX , byteimg.com  ) )  ")
	t.Log("re", re)

}

func TestRule2(t *testing.T) {
	re := CheckURLREGEX("http://google1.com", `^http://google\.com`)
	t.Log(re)

}

func TestJq(t *testing.T) {
	//json := `{"name": "John", "age": 30, "city": "New York"}`
	jsonmap := map[string]any{"name": "John", "age": 30, "city": "New York"}
	//query := ".name"
	query, err := gojq.Parse("del(.name)")
	if err != nil {
		t.Log("err", err)
		log.Fatalln(err)
	}
	//input := map[string]any{"foo": []any{1, 2, 3}}
	iter := query.Run(jsonmap) // or query.RunWithContext
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
				break
			}
			t.Log("err2", err)
			log.Fatalln(err)
		}
		fmt.Printf("%#v\n", v)
		t.Log("v", v)
	}
}
