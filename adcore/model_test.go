package adcore

import (
	"fmt"
	"log"
	"net/url"
	"testing"

	"github.com/itchyny/gojq"
)

func TestP12(t *testing.T) {

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
