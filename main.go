package main

import (
	"context"
	"github.com/Kumengda/laser/crawler"
	. "github.com/Kumengda/laser/runtime"
)

func main() {
	myCrawler, _ := crawler.NewCrawler("https://news.sina.com.cn",
		2,
		[]string{},
	)
	//chromeCrawler, _ := crawler.NewChromeCrawler(false, 5, false, 30, 10, nil)
	nativeCrawler := crawler.NewNativeCrawler(10, 10, nil)
	myCrawler.SetCrawler(nativeCrawler)
	myCrawler.SetMiddlewareFunc(func(i interface{}) interface{} {
		return i
	})
	res := myCrawler.ParamCrawl(context.Background())
	MainInsp.Print(Json(res))
	//fmt.Println(len(res.ExternalLink) + len(res.ExternalStaticFileLink) + len(res.SameOriginUrl))
}
