package crawler

import (
	"github.com/Kumengda/easyChromedp/chrome"
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
	SameOriginForm         []template.JsRes         `json:"same_origin_form"`
	ExternalForm           []template.JsRes         `json:"external_form"`
}

type ChromeCrawler struct {
	headers       map[string]interface{}
	waitTime      int
	printLog      bool
	timeout       int
	chromeThreads int
	chrome        *chrome.Chrome
}

func NewChromeCrawler(printLog bool, waitTime int, headless bool, chromeThreads, timeout int, headers map[string]interface{}) (*ChromeCrawler, error) {
	if headers == nil {
		headers = make(map[string]interface{})
	}
	myChrome, err := chrome.NewChrome(
		chromedp.Flag("headless", headless),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
	)

	if err != nil {
		return nil, err
	}
	return &ChromeCrawler{
		headers:       headers,
		timeout:       timeout,
		waitTime:      waitTime,
		printLog:      printLog,
		chromeThreads: chromeThreads,
		chrome:        myChrome,
	}, nil
}
func (c *ChromeCrawler) GetCrawlThreads() int {
	return c.chromeThreads
}
func (c *ChromeCrawler) Close() {
	c.chrome.Close()
}
func (c *ChromeCrawler) DoFinally() {
	c.chrome.Close()
}

func (c *ChromeCrawler) SingleCrawl(task template.JsRes, allHref []template.JsRes) []template.JsRes {
	templates, err := template.NewChromedpTemplates(
		c.printLog,
		c.timeout,
		c.waitTime,
		c.headers,
		c.chrome,
	)
	if err != nil {
		return allHref
	}
	allReqHref, err := templates.GetWebsiteAllReq(task.Url)
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
	allJsHref, err := templates.GetWebsiteAllHrefByJs(task.Url)
	if err != nil {
		return nil
	}
	allHref = append(allHref, allJsHref...)
	return allHref
}
