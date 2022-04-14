package core

import (
	"fmt"
	"net/url"
	"regexp"
	"time"
)

type TestResult struct {
	Id       string
	url      url.URL
	status   string
	duration time.Duration
}

type Tester struct {
	requestTimeout time.Duration
	testInterval time.Duration
	out chan TestResult
}

func NewTester(config WatchConfig, out chan TestResult) Tester {
	return Tester { 
		requestTimeout: time.Duration(config.RequestTimeout) * time.Second, 
		testInterval: time.Duration(config.TestInterval) * time.Second,
		out: out,
	}
}

func (t *Tester)Test(url *url.URL) {
	tp := NewTransport(t.requestTimeout)

	for {
		t.out <- TestResult{Id: url.Host, url: *url, status: "Test"}

		_, err := tp.Dial(url.Scheme, url.Host)

		if err != nil {
			t.out <- TestResult{Id: url.Host, url: *url, status: formatError(err, url), duration: tp.ConnDuration()}
			return
		}

		t.out <- TestResult{Id: url.Host, url: *url, status: "OK", duration: tp.Duration()}
		time.Sleep(t.testInterval)

	}
}

func formatError(err error, url *url.URL) string {
	m := regexp.MustCompile(fmt.Sprintf(`(net/http: )| \(.*\)|(dial .* %s: )`, url.Host))
	return m.ReplaceAllString(err.Error(), "")
}
