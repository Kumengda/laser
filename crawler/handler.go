package crawler

import (
	"github.com/Kumengda/easyChromedp/template"
	. "github.com/Kumengda/laser/runtime"
	"regexp"
	"sort"
	"strings"
)

var noCrawlExtension = []string{".png",
	".css",
	".js",
	".jpg",
	".png",
	".jpeg",
	".gif",
	".svg",
	".ttf",
	".otf",
	".woff",
	".woff2",
	".mp4",
	".avi",
	".mov",
	".wmv",
	".mp3",
	".wav",
	".ogg",
	".csv",
	".json",
	".xml",
	".xls",
	".xlsx",
}

func splitStringSlice(slice []string, x int) [][]string {
	n := len(slice)
	if n < x {
		x = n
	}

	result := make([][]string, x)
	chunkSize := n / x // 计算每份的大小
	remainder := n % x // 切片剩余部分的长度

	index := 0
	for i := 0; i < x; i++ {
		chunk := chunkSize
		if i < remainder {
			chunk++
		}
		result[i] = slice[index : index+chunk]
		index += chunk
	}

	return result
}

func compareCompareJsRes(a, b []template.JsRes) bool {
	if len(a) != len(b) {
		return false
	}
	var a1 []string
	var b1 []string
	for _, v := range a {
		a1 = append(a1, v.Url)
	}
	for _, v := range b {
		b1 = append(b1, v.Url)
	}
	sort.Strings(a1)
	sort.Strings(b1)

	for i := range a1 {
		if a1[i] != b1[i] {
			return false
		}
	}

	return true
}
func containsString(list []template.JsRes, target string) bool {
	for _, str := range list {
		if str.Url == target {
			return true
		}
	}
	return false
}

func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
func mapsEqual(a, b []template.FormData) bool {
	if len(a) != len(b) {
		return false
	}
	var check bool
	for _, v := range a {
		check = false
		for _, v1 := range b {
			if v1.Name == v.Name {
				check = true
				break
			}
		}
		if !check {
			return false
		}
	}
	return check
}

func staticCheck(target string) bool {
	for _, e := range noCrawlExtension {
		if strings.HasSuffix(target, e) {
			return true
		}
	}
	return false
}

func continueCheck(target string, host string, filter []string) bool {
	for _, v := range filter {
		re := regexp.MustCompile(v)
		if re.MatchString(target) {
			MainInsp.Print(LEVEL_INFO, Text("黑名单命中:"+target))
			return false
		}
	}
	for _, v := range noCrawlExtension {
		if strings.HasSuffix(target, v) {
			return false
		}
	}
	thost, err := getHost(target)
	if err != nil {
		return false
	}
	if thost != host {
		return false
	}
	return true
}
func targetRemoveDuplicates(target []template.JsRes) []template.JsRes {
	var res []template.JsRes
	for _, v1 := range target {
		var check bool
		for _, v2 := range res {
			if v2.Url == v1.Url && v2.IsForm == v1.IsForm && v2.Method == v1.Method && mapsEqual(v2.Param, v1.Param) {
				check = true
			}
		}
		if !check {
			res = append(res, v1)
		}
	}
	return res
}
