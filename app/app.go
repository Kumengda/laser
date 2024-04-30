package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/B9O2/tabby"
	"github.com/Kumengda/easyChromedp/template"
	"github.com/Kumengda/laser/crawler"
	"github.com/Kumengda/laser/mySignal"
	. "github.com/Kumengda/laser/runtime"
	"os"
	"os/signal"
	"syscall"
)

type MainApp struct {
	*tabby.BaseApplication
	version string
	sigChan chan os.Signal
	tips    string
}

func NewMainApp(version, tips string) *MainApp {
	return &MainApp{BaseApplication: tabby.NewBaseApplication(0, 0, nil),
		version: version,
		tips:    tips,
		sigChan: make(chan os.Signal),
	}
}
func (m MainApp) Detail() (string, string) {
	return "scanner", "detail"
}
func (m *MainApp) Down() {
	m.sigChan <- mySignal.NewAppDownSignal()
}
func (m MainApp) Main(arguments tabby.Arguments) (*tabby.TabbyContainer, error) {
	if arguments.IsEmpty() {
		m.Help(m.version + "\n" + m.tips)
		return nil, nil
	}
	target := arguments.Get("target").(string)
	depth := arguments.Get("depth").(int)
	filter := arguments.Get("noCrawler").([]string)
	headers := arguments.Get("headers").(map[string]interface{})
	enableChromeCrawler := arguments.Get("enableChromeCrawler").(bool)
	timeout := arguments.Get("timeout").(int)
	threads := arguments.Get("threads").(int)
	filename := arguments.Get("file").(string)
	headless := arguments.Get("headless").(bool)
	waitTime := arguments.Get("waitTime").(int)
	myCrawler, _ := crawler.NewCrawler(target,
		depth,
		filter,
	)
	var chromeCrawler *crawler.ChromeCrawler
	var nativeCrawler *crawler.NativeCrawler
	var err error
	go func() {
		signal.Notify(m.sigChan, os.Interrupt, syscall.SIGTERM)
		sig := <-m.sigChan
		switch sig {
		case os.Interrupt:
			MainInsp.Print(LEVEL_WARNING, Text("程序正在尝试停止,请等待......"))
			if chromeCrawler != nil {
				chromeCrawler.Close()
				os.Exit(0)
			}
		case mySignal.AppDownSignal{}:
			return
		}
	}()

	if enableChromeCrawler {
		chromeCrawler, err = crawler.NewChromeCrawler(false, waitTime, headless, threads, timeout, headers)
		if err != nil {
			return nil, err
		}
		myCrawler.SetCrawler(chromeCrawler)
	} else {
		nativeCrawler = crawler.NewNativeCrawler(timeout, threads, headers)
		myCrawler.SetCrawler(nativeCrawler)
	}
	myCrawler.SetMiddlewareFunc(func(i interface{}) interface{} {
		if i != nil {
			res := i.([]template.JsRes)
			for _, v := range res {
				if v.Method == "" {
					v.Method = "GET"
				}
				MainInsp.Print(LEVEL_INFO, Text(fmt.Sprintf("[%s]%s IsForm:%t", v.Method, v.Url, v.IsForm)))
			}
		}
		return i
	})
	res := myCrawler.ParamCrawl(context.Background())
	if filename != "" {
		jsonData, _ := json.Marshal(res)
		os.WriteFile(filename, jsonData, 0644)
	}
	m.Down()
	return nil, nil
}
