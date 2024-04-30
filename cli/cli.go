package main

import (
	"github.com/B9O2/tabby"
	"github.com/Kumengda/laser/app"
	. "github.com/Kumengda/laser/runtime"
	"github.com/tidwall/gjson"
)

var RegExp = tabby.NewTransfer("regexp", func(s string) (any, error) {
	var regexp = []string{}
	data := gjson.Parse(s)
	for _, v := range data.Array() {
		regexp = append(regexp, v.String())
	}
	return regexp, nil
})

var Headers = tabby.NewTransfer("headers", func(s string) (any, error) {
	var headers = make(map[string]interface{})
	data := gjson.Parse(s)
	for k, v := range data.Map() {
		headers[k] = v.String()
	}
	return headers, nil
})

func main() {
	mainApp := app.NewMainApp("v0.0.1", "*表示必要参数")
	mainApp.SetParam("target", "*目标,使用,分割多个扫描目标", tabby.String(nil), "tg")
	mainApp.SetParam("enableChromeCrawler", "爬取深度,默认2", tabby.Bool(false), "ec")
	mainApp.SetParam("headless", "是否启用headless模式 默认false", tabby.Bool(false), "hl")
	mainApp.SetParam("depth", "爬取深度,默认2", tabby.Int(2), "d")
	mainApp.SetParam("timeout", "超时时间,默认5", tabby.Int(5), "ti")
	mainApp.SetParam("waitTime", "chrome等待js加载时间,默认2", tabby.Int(2), "wt")
	mainApp.SetParam("threads", "爬虫线程,默认20", tabby.Int(20), "th")
	mainApp.SetParam("noCrawler", "不进行爬取的url规则,采用正则形式,入参规则eg:[\"regexp1\",\"regexp2\"]", RegExp([]string{}), "nc")
	mainApp.SetParam("headers", "请求头,请使用json格式,eg:{\"key\":\"value\"}", Headers(map[string]interface{}{}), "hs")
	mainApp.SetParam("file", "结果保存文件路径", tabby.String(""), "file")
	t := tabby.NewTabby("laser", mainApp)
	_, err := t.Run(nil)
	if err != nil {
		MainInsp.Print()
		return
	}
}
