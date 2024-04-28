package crawler

import (
	"context"
	"github.com/B9O2/Multitasking"
	"github.com/Kumengda/easyChromedp/template"
	. "github.com/Kumengda/laser/runtime"
	"strings"
)

type Crawler struct {
	BaseCrawler
	crawl BaseCrawl
}

func NewCrawler(target string, depth int, noCrawlerFilter []string) (*Crawler, error) {
	Init()
	host, err := getHost(target)
	if err != nil {
		return nil, err
	}
	return &Crawler{BaseCrawler: BaseCrawler{
		target:         target,
		depth:          depth,
		host:           host,
		filter:         noCrawlerFilter,
		middlewareFunc: nil,
	}}, nil
}
func (c *Crawler) SetCrawler(crawler BaseCrawl) {
	c.crawl = crawler
}

func (c *Crawler) Crawl() DirResult {
	var dirRes DirResult
	dirRes.Target = c.target
	res := c.crawAllUrl([]template.JsRes{{Url: c.target, Method: "GET"}}, nil, context.Background(), c.crawl)
	myMT := Multitasking.NewMultitasking("urlCheck", nil)
	myMT.Register(func(dc Multitasking.DistributeController) {
		for _, v := range res {
			parse, err := getHost(v.Url)
			if err != nil {
				continue
			}
			if parse == c.host {
				dc.AddTask(SameOriginUrl{BaseUrl{Url: v.Url}})
			} else {
				if staticCheck(v.Url) {
					dc.AddTask(ExternalStaticFileLink{BaseUrl: BaseUrl{Url: v.Url}})
					continue
				}
				dc.AddTask(ExternalLink{BaseUrl: BaseUrl{Url: v.Url}})
			}
		}
	}, func(ec Multitasking.ExecuteController, i interface{}) interface{} {
		if c.middlewareFunc != nil {
			midres := c.middlewareFunc(i)
			if midres != nil {
				return midres
			}
		}
		return i
	})

	runRes, err := myMT.Run(context.Background(), 50)
	if err != nil {
		return dirRes
	}
	for _, v := range runRes {
		if v != nil {
			switch v.(type) {
			case SameOriginUrl:
				dirRes.SameOriginUrl = append(dirRes.SameOriginUrl, v.(SameOriginUrl))
			case ExternalLink:
				dirRes.ExternalLink = append(dirRes.ExternalLink, v.(ExternalLink))
			case ExternalStaticFileLink:
				dirRes.ExternalStaticFileLink = append(dirRes.ExternalStaticFileLink, v.(ExternalStaticFileLink))
			}

		}
	}
	return dirRes
}
func (c *Crawler) ParamCrawl() []template.JsRes {
	var sameOriginRes []template.JsRes
	res := c.crawAllUrl([]template.JsRes{{Url: c.target, Method: "GET"}}, nil, context.Background(), c.crawl)
	for _, v := range res {
		switch v.Method {
		case "GET":
			if v.IsForm && len(v.Param) == 0 {
				continue
			}
			if !v.IsForm && !strings.Contains(v.Url, "?") {
				continue
			}
		case "POST":
			if len(v.Param) == 0 {
				continue
			}
		}
		parse, err := getHost(v.Url)
		if err != nil {
			continue
		}
		if parse == c.host {
			sameOriginRes = append(sameOriginRes, v)
		}
	}
	return sameOriginRes
}
