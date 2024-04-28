package crawler

import (
	"github.com/Kumengda/easyChromedp/template"
	"github.com/Kumengda/laser/http"
	"github.com/Kumengda/pageParser/parser"
	"net/url"
	"strings"
)

type NativeCrawler struct {
	headers map[string]interface{}
	timeout int
	threads int
}

func (n *NativeCrawler) SingleCrawl(task template.JsRes, allHref []template.JsRes) []template.JsRes {
	resp, err, _ := http.Get(task.Url, n.headers, n.timeout)
	if err != nil {
		return nil
	}
	parse, err := url.Parse(task.Url)
	if err != nil {
		return nil
	}
	host := parse.Host
	scheme := parse.Scheme
	tagExtract := parser.NewTagExtract()
	tagExtract.InitTags(parser.DefaultTagRules)
	fromExtract := parser.NewFormExtract()
	tagRes := tagExtract.Extract(string(resp))
	formRes := fromExtract.Extract(string(resp))
	var allTagUrl []string
	for _, t := range tagRes {
		for _, v := range t.Attr {
			allTagUrl = append(allTagUrl, v)
		}
	}
	allTagUrl = cleanUrl(allTagUrl)
	allTagUrl = removeDuplicateStrings(allTagUrl)
	for _, v := range allTagUrl {
		allHref = append(allHref, template.JsRes{
			Url:    parseHrefData(v, scheme, host, task.Url, false),
			Method: "GET",
		})
	}
	for _, v := range formRes {
		var fromUrl string
		var newFormData []template.FormData
		isFileUpload := false
		for _, vv := range v.FormData {
			if vv.Name == "" || !checkInputType(vv.Type) {
				continue
			}
			if vv.Type == "file" {
				isFileUpload = true
			}
			newFormData = append(newFormData, template.FormData{
				Name:  vv.Name,
				Type:  vv.Type,
				Value: vv.Value,
			})
		}
		if v.Action == "#" || v.Action == "/" || v.Action == "" {
			fromUrl = task.Url
		} else {
			fromUrl = parseHrefData(v.Action, scheme, host, task.Url, true)
		}
		allHref = append(allHref, template.JsRes{
			Url:          fromUrl,
			Method:       strings.ToUpper(v.Method),
			IsForm:       true,
			Param:        newFormData,
			IsFileUpload: isFileUpload,
		})
	}
	return allHref
}

func (n *NativeCrawler) GetCrawlThreads() int {
	return n.threads
}

func NewNativeCrawler(timeout, threads int, headers map[string]interface{}) *NativeCrawler {
	if headers == nil {
		headers = make(map[string]interface{})
	}
	return &NativeCrawler{
		timeout: timeout,
		threads: threads,
		headers: headers,
	}
}
