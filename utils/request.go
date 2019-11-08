package utils

import (
	"golang.org/x/net/proxy"
	"hammer-x/config"
	"net"
	"net/http"
	"net/url"
	netURL "net/url"
	"strings"
	"time"
)

// Request base request
func Curl(
	method, url string, data string, headers map[string]string,
) (*http.Response, error) {
	transport := &http.Transport{
		DisableCompression: true,
	}
	if config.Proxy != "" {
		var httpProxy, err = netURL.Parse(config.Proxy)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(httpProxy)
	}
	if config.Socket != "" {
		dialer, err := proxy.SOCKS5(
			"tcp",
			config.Socket,
			nil,
			&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			},
		)
		if err != nil {
			return nil, err
		}
		transport.Dial = dialer.Dial
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Minute,
	}
	var req *http.Request
	var err error
	if data == "" {
		urlArr := strings.Split(url, "?")
		if len(urlArr) == 2 {
			url = urlArr[0] + "?" + getParseParam(urlArr[1])
		}
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, strings.NewReader(data))
	}
	if err != nil {
		return nil, err
	}
	for k, v := range config.Headers {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if _, ok := headers["Referer"]; !ok {
		req.Header.Set("Referer", url)
	}

	res, err := client.Do(req)
	defer res.Body.Close()
	return res, err

	//if config.Cookie != "" {
	//	var cookie string
	//	var cookies []*http.Cookie
	//	cookies, err = cookiemonster.ParseString(config.Cookie)
	//	if err != nil || len(cookies) == 0 {
	//		cookie = config.Cookie
	//	}
	//	if cookie != "" {
	//		req.Header.Set("Cookie", cookie)
	//	}
	//	if cookies != nil {
	//		for _, c := range cookies {
	//			req.AddCookie(c)
	//		}
	//	}
	//}
	//if config.Refer != "" {
	//	req.Header.Set("Referer", config.Refer)
	//}

}

func getParseParam(param string) string {
	return url.PathEscape(param)
}
