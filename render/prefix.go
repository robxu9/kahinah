package render

import (
	"bytes"
	"html/template"
	"net/url"
	"strings"

	"github.com/robxu9/kahinah/common/conf"
	"github.com/robxu9/kahinah/common/klog"
)

var (
	baseURL = conf.Get("baseURL").(string)
)

// ConvertURLWithData merges a relative URL into an absolute one,
// also templating as needed.
func ConvertURLWithData(dest string, data interface{}) string {
	// no need to prefix if the dest has no / before it
	temp := template.Must(template.New("prefixTemplate").Parse(dest))
	var b bytes.Buffer

	err := temp.Execute(&b, data)
	if err != nil {
		panic(err)
	}

	result := b.String()
	return ConvertURL(result)
}

// ConvertURL converts a URL from a relative into an absolute location.
func ConvertURL(dest string) string {
	baseTrim := strings.TrimRight(baseURL, "/")
	destTrim := strings.TrimLeft(dest, "/")
	return baseTrim + "/" + destTrim
}

// ConvertURLRelative converts a URL to a better relative location (i.e. one
// that works without the host but keeps the trailing path).
func ConvertURLRelative(dest string) string {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		klog.Panicf("couldn't parse the baseurl - is it valid??: %v", err)
	}

	if strings.Trim(parsedURL.RawPath, "/") != "" {
		return "/" + strings.Trim(parsedURL.RawPath, "/") + "/" + strings.TrimLeft(dest, "/")
	}

	return "/" + strings.TrimLeft(dest, "/")
}
