package core

import (
	"fmt"
	"net/url"
	"time"
)

type TestResult struct {
	Id   string
	url  url.URL
	ping string
}

type Tester struct {
	requestTimeout time.Duration
	testInterval   time.Duration
	out            chan TestResult
	quit           chan bool
}

func NewTester(config *WatchConfig, out chan TestResult, quit chan bool) Tester {
	return Tester{
		requestTimeout: time.Duration(config.RequestTimeout) * time.Second,
		testInterval:   time.Duration(config.TestInterval) * time.Second,
		out:            out,
		quit:           quit,
	}
}

func (t *Tester) Test(url *url.URL) {
	for {
		select {
		case <-t.quit:
			return
		default:
			pass := t.ping(url)
			t.out <- TestResult{Id: url.Host, url: *url, ping: fmt.Sprintf("%d/10", pass)}
			time.Sleep(t.testInterval)
		}

	}
}

func (t *Tester) ping(url *url.URL) int {
	tp := NewTransport(t.requestTimeout)
	pass := 0
	for i := 0; i < 10; i++ {
		_, err := tp.Dial(url.Scheme, url.Host)
		if err == nil {
			pass++
		}
	}

	return pass
}
