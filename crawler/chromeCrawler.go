package crawler

import (
	"github.com/Kumengda/easyChromedp/template"
	"github.com/chromedp/chromedp"
)

type BaseUrl struct {
	Url string `json:"url"`
}

type SameOriginUrl struct {
	BaseUrl
}

type ExternalLink struct {
	BaseUrl
}

type ExternalStaticFileLink struct {
	BaseUrl
}

type DirResult struct {
	Target                 string                   `json:"target"`
	SameOriginUrl          []SameOriginUrl          `json:"same_originUrl"`
	ExternalLink           []ExternalLink           `json:"external_link"`
	ExternalStaticFileLink []ExternalStaticFileLink `json:"external_static_file_link"`
}

type ChromeCrawler struct {
	timeout       int
	headers       map[string]interface{}
	waitTime      int
	printLog      bool
	chromeThreads int
	headless      bool
}

func NewChromeCrawler(printLog bool, waitTime int, headless bool, chromeThreads, timeout int, headers map[string]interface{}) (*ChromeCrawler, error) {
	if headers == nil {
		headers = make(map[string]interface{})
	}
	return &ChromeCrawler{
		timeout:       timeout,
		headers:       headers,
		waitTime:      waitTime,
		printLog:      printLog,
		chromeThreads: chromeThreads,
		headless:      headless,
	}, nil
}
func (c *ChromeCrawler) GetCrawlThreads() int {
	return c.chromeThreads
}
func (c *ChromeCrawler) SingleCrawl(task template.JsRes, allHref []template.JsRes) []template.JsRes {
	templates, err := template.NewChromedpTemplates(
		task.Url,
		c.timeout,
		c.printLog,
		c.waitTime,
		c.headers,
		chromedp.Flag("headless", c.headless),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
	)
	if err != nil {
		return allHref
	}
	allReqHref, err := templates.GetWebsiteAllReq()
	if err != nil {
		return nil
	}
	for _, v := range allReqHref {
		allHref = append(allHref, template.JsRes{
			Url:    v,
			Method: "GET",
			Param:  nil,
			IsForm: false,
		})

	}
	allJsHref, err := templates.GetWebsiteAllHrefByJs()
	if err != nil {
		return nil
	}
	allHref = append(allHref, allJsHref...)
	return allHref
}
