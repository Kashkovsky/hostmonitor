package core

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
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

type TestResult struct {
	url    url.URL
	status string
}

func Test(url *url.URL, timeoutSeconds int, out chan TestResult) {
	timeOut := time.Duration(timeoutSeconds) * time.Second
	_, err := net.DialTimeout(url.Scheme, url.Host, timeOut)

	if err != nil {
		out <- TestResult{url: *url, status: err.Error()}
		return
	}

	out <- TestResult{url: *url, status: "OK"}
}
