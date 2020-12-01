package em_utils

import (
	"net/url"
	"strings"
)

func GetUrlPath(urlStr string, trim bool) string {
	u, _ :=url.Parse(urlStr)
	if trim {
		return strings.TrimLeft(u.Path, "/")
	} else {
		return u.Path
	}
}