package utils

import (
	"compress/flate"
	"compress/gzip"
	"errors"
	"fmt"
	"golang.org/x/net/proxy"
	"hammer-x/config"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	netURL "net/url"
	"strconv"
	"strings"
	"time"
)

// Request base request
func Request(
	method, url string, body io.Reader, headers map[string]string,
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
	req, err := http.NewRequest(method, url, body)
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

	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		err = fmt.Errorf("request error: %v", err)
	}
	return res, nil
}

// Get get request
func Get(url, refer string, headers map[string]string) (string, error) {
	body, err := GetByte(url, refer, headers)
	return string(body), err
}

// GetByte get request
func GetByte(url, refer string, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	if refer != "" {
		headers["Referer"] = refer
	}
	res, err := Request(http.MethodGet, url, nil, headers)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(res.Body)
	case "deflate":
		reader = flate.NewReader(res.Body)
	default:
		reader = res.Body
	}
	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return body, nil
}

// Headers return the HTTP Headers of the url
func Headers(url, refer string) (http.Header, error) {
	headers := map[string]string{
		"Referer": refer,
	}
	res, err := Request(http.MethodGet, url, nil, headers)
	if err != nil {
		return nil, err
	}
	return res.Header, nil
}

// Size get size of the url
func Size(url, refer string) (int64, error) {
	h, err := Headers(url, refer)
	if err != nil {
		return 0, err
	}
	s := h.Get("Content-Length")
	if s == "" {
		return 0, errors.New("Content-Length is not present")
	}
	size, _ := strconv.Atoi(s)
	downloadSize := int64(size)
	return downloadSize, nil
}

// ContentType get Content-Type of the url
func ContentType(url, refer string) (string, error) {
	h, err := Headers(url, refer)
	if err != nil {
		return "", err
	}
	s := h.Get("Content-Type")
	// handle Content-Type like this: "text/html; charset=utf-8"
	return strings.Split(s, ";")[0], nil
}
