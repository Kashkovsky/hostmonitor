package core

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
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
	Id       string
	url      url.URL
	status   string
	duration time.Duration
}

func Test(url *url.URL, timeoutSeconds int, testInterval int, out chan TestResult) {
	timeout := time.Duration(timeoutSeconds) * time.Second
	tp := NewTransport(timeout)

	for {
		out <- TestResult{Id: url.Host, url: *url, status: "Test"}

		_, err := tp.Dial(url.Scheme, url.Host)

		if err != nil {
			out <- TestResult{Id: url.Host, url: *url, status: formatError(err, url), duration: tp.ConnDuration()}
			return
		}

		out <- TestResult{Id: url.Host, url: *url, status: "OK", duration: tp.Duration()}
		time.Sleep(time.Duration(testInterval) * time.Second)

	}
}

func formatError(err error, url *url.URL) string {
	m := regexp.MustCompile(fmt.Sprintf(`(net/http: )| \(.*\)|(dial .* %s: )`, url.Host))
	return m.ReplaceAllString(err.Error(), "")
}
