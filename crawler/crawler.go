package crawler

import (
	"context"
	"github.com/Kumengda/easyChromedp/template"
	. "github.com/Kumengda/laser/runtime"
	"strings"
)

type Crawler struct {
	BaseCrawler
	crawl Crawl
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
func (c *Crawler) SetCrawler(crawler Crawl) {
	c.crawl = crawler
}

func (c *Crawler) Crawl() DirResult {
	var dirRes DirResult
	dirRes.Target = c.target
	res := c.crawAllUrl([]template.JsRes{{Url: c.target, Method: "GET"}}, nil, context.Background(), c.crawl)
	for _, v := range res {
		parse, err := getHost(v.Url)
		if err != nil {
			continue
		}
		if v.IsForm {
			if parse == c.host {
				dirRes.SameOriginForm = append(dirRes.SameOriginForm, v)
			} else {
				dirRes.ExternalForm = append(dirRes.ExternalForm, v)
			}
			continue
		}

		if parse == c.host {
			dirRes.SameOriginUrl = append(dirRes.SameOriginUrl, SameOriginUrl{BaseUrl{Url: v.Url}})
		} else {
			if staticCheck(v.Url) {
				dirRes.ExternalStaticFileLink = append(dirRes.ExternalStaticFileLink, ExternalStaticFileLink{BaseUrl: BaseUrl{Url: v.Url}})
				continue
			}
			dirRes.ExternalLink = append(dirRes.ExternalLink, ExternalLink{BaseUrl: BaseUrl{Url: v.Url}})
		}
	}
	c.crawl.DoFinally()
	return dirRes
}
func (c *Crawler) ParamCrawl(ctx context.Context) []template.JsRes {
	var sameOriginRes []template.JsRes
	res := c.crawAllUrl([]template.JsRes{{Url: c.target, Method: "GET"}}, nil, ctx, c.crawl)
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
	c.crawl.DoFinally()
	return sameOriginRes
}
