package main

import (
	"context"
	"github.com/Kumengda/laser/crawler"
	. "github.com/Kumengda/laser/runtime"
)

func main() {
	myCrawler, _ := crawler.NewCrawler("http://127.0.0.1:8765/",
		2,
		[]string{".*logout.*", ".*lang.*"},
	)
	chromeCrawler, _ := crawler.NewChromeCrawler(false, 5, true, 10, 10, nil)
	//nativeCrawler := crawler.NewNativeCrawler(10, 10, nil)
	myCrawler.SetCrawler(chromeCrawler)
	myCrawler.SetMiddlewareFunc(func(i interface{}) interface{} {
		MainInsp.Print(Json(i))
		return i
	})
	myCrawler.ParamCrawl(context.Background())

	//fmt.Println(len(res.ExternalLink) + len(res.ExternalStaticFileLink) + len(res.SameOriginUrl))
}
