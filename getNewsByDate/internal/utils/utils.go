package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// GetPage call the client page by HTTP request and extract the body to HTML document.
func GetPage(
	ctx context.Context,
	method,
	siteURL string,
	cookies []*http.Cookie,
	headers,
	formDatas map[string]string,
	timeout int,
) (*goquery.Document, []*http.Cookie, error) {

	body := io.Reader(nil)
	if len(formDatas) > 0 {
		form := url.Values{}
		for k, v := range formDatas {
			form.Add(k, v)
		}

		body = strings.NewReader(form.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, method, siteURL, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create http request context: %w", err)
	}

	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	if len(cookies) > 0 {
		for _, c := range cookies {
			req.AddCookie(c)
		}
	}

	reqTimeout := 10 * time.Second
	if timeout != 0 {
		reqTimeout = time.Duration(timeout) * time.Second
	}

	httpClient := &http.Client{
		Transport:     http.DefaultTransport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       reqTimeout,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute http request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse html: %w", err)
	}
	return doc, resp.Cookies(), nil
}
