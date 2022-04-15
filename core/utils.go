package core

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetStringFromURL(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func ParceUrls(urls string) []*url.URL {
	records := strings.Split(urls, "\n")
	results := []*url.URL{}
	for _, addr := range records {
		res, err := parceUrl(addr)
		if err != nil {
			continue
		}

		results = append(results, res)
	}

	return results
}

func parceUrl(addr string) (*url.URL, error) {
	if addr == "" {
		return nil, errors.New("URL is an empty string")
	}
	u, err := url.Parse(addr)
	if err != nil {
		if !strings.Contains(addr, "://") {
			return parceUrl("tcp://" + addr)
		}
		return nil, err
	}

	return u, nil
}
