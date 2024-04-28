package crawler

import (
	"context"
	"fmt"
	"github.com/B9O2/Multitasking"
	"github.com/Kumengda/easyChromedp/template"
	. "github.com/Kumengda/laser/runtime"
)

type BaseCrawl interface {
	SingleCrawl(task template.JsRes, allHref []template.JsRes) []template.JsRes
	GetCrawlThreads() int
}

type BaseCrawler struct {
	target         string
	timeout        int
	depth          int
	host           string
	filter         []string
	middlewareFunc func(i interface{}) interface{}
}

func (b *BaseCrawler) SetMiddlewareFunc(middlewareFunc func(i interface{}) interface{}) {
	b.middlewareFunc = middlewareFunc
}

// 第一次入参targets只能是一个包含一个目标的切片类型,这里第一个参数写成[]string是方便递归传参
func (b *BaseCrawler) crawAllUrl(targets []template.JsRes, lastTargets []template.JsRes, ctx context.Context, crawl BaseCrawl) []template.JsRes {
	targets = targetRemoveDuplicates(targets)
	var allHref []template.JsRes
	if lastTargets != nil {
		allHref = append(allHref, targets...)
	}
	if b.depth == 0 {
		return targets
	}
	b.depth = b.depth - 1
	if compareCompareJsRes(targets, lastTargets) {
		return targets
	}
	myMT := Multitasking.NewMultitasking("crawler", nil)
	myMT.Register(func(dc Multitasking.DistributeController) {
		for _, v := range targets {
			if lastTargets == nil {
				//证明是第一次
				dc.AddTask(v)
			} else {
				if !containsString(lastTargets, v.Url) {
					dc.AddTask(v)
				}
			}
		}
	}, func(ec Multitasking.ExecuteController, i interface{}) interface{} {

		select {
		case <-ctx.Done():
			return nil
		default:

		}
		var _allHref []template.JsRes
		task := i.(template.JsRes)
		if !continueCheck(task.Url, b.host, b.filter) {
			return append(_allHref, template.JsRes{Url: task.Url, Method: "GET"})
		}
		MainInsp.Print(LEVEL_INFO, Text(fmt.Sprintf("depth:%d Req:%s", b.depth, task.Url)))
		return crawl.SingleCrawl(task, _allHref)
	})
	myMT.SetResultMiddlewares(Multitasking.NewBaseMiddleware(func(ec Multitasking.ExecuteController, i interface{}) (interface{}, error) {
		if b.middlewareFunc != nil {
			midres := b.middlewareFunc(i)
			if midres == nil {
				return i, nil
			}
			return midres, nil
		}
		return i, nil
	}))
	res, err := myMT.Run(context.Background(), uint(crawl.GetCrawlThreads()))
	if err != nil {
		return nil
	}
	for _, v := range res {
		if v != nil {
			allHref = append(allHref, v.([]template.JsRes)...)
		}
	}
	return b.crawAllUrl(allHref, targets, ctx, crawl)
}
