package crawler

import (
	"net/url"
	"strings"
)

var typeList = []string{
	"button",
	"checkbox",
	"color",
	"radio",
	"range",
	"reset",
	"search",
}

func getHost(target string) (string, error) {
	parse, err := url.Parse(target)
	if err != nil {
		return "", err
	}
	host := parse.Host
	if strings.Contains(host, ":443") {
		host = strings.Split(host, ":")[0]
	}
	if strings.Contains(target, ":80") {
		host = strings.Split(host, ":")[0]
	}
	host = parse.Scheme + "://" + host
	return host, nil
}
func cleanUrl(target []string) []string {
	var newUrls = []string{}
	for _, v := range target {
		if replace(v, "", "#", "\n", " ") != "" {
			newUrls = append(newUrls, v)
		}
	}
	return newUrls
}
func replace(rawStr string, replaceStr ...string) string {
	for _, v := range replaceStr {
		rawStr = strings.ReplaceAll(rawStr, v, "")
	}
	return rawStr
}
func parseHrefData(u, scheme, host, nowUrl string, isForm bool) string {

	//if strings.Contains(u, "?") {
	//	u = u[:strings.Index(u, "?")]
	//}
	if strings.HasPrefix(u, "//") {
		return scheme + ":" + u
	}
	if strings.HasPrefix(u, "/") {
		return scheme + "://" + host + u
	}
	if strings.HasPrefix(u, "#") {
		return nowUrl + u

	}
	if strings.HasPrefix(u, "http") {
		return u
	}
	if isForm && strings.Contains(nowUrl, "#") {
		nowUrl = nowUrl[:strings.Index(nowUrl, "#")]
	}
	nowUrl = nowUrl[:strings.LastIndex(nowUrl, "/")]
	return nowUrl + "/" + u
}
func removeDuplicateStrings(strList []string) []string {
	stringMap := make(map[string]bool)
	for _, str := range strList {
		stringMap[str] = true
	}
	newStrList := []string{}
	for key := range stringMap {
		if key != "" {
			newStrList = append(newStrList, key)
		}
	}
	return newStrList
}
func checkInputType(t string) bool {
	check := strings.ToLower(t)
	for _, v := range typeList {
		if check == v {
			return false
		}
	}
	return true
}
