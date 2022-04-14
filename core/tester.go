package core

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type TestResult struct {
	Id         string
	InProgress bool
	url        url.URL
	tcp        string
	httpStatus string
	duration   string
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
			t.out <- TestResult{
				Id:         url.Host,
				InProgress: true,
				httpStatus: "Testing...",
			}
			pass := t.tcp(url)
			status, duration := t.http(url)
			t.out <- TestResult{
				Id:         url.Host,
				url:        *url,
				tcp:        fmt.Sprintf("%d/10", pass),
				httpStatus: status,
				duration:   strconv.FormatInt(duration.Milliseconds(), 10) + "ms",
			}
		}
		time.Sleep(t.testInterval)
	}
}

func (t *Tester) tcp(url *url.URL) int {
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

func (t *Tester) http(url *url.URL) (status string, duration time.Duration) {
	tp := NewTransport(t.requestTimeout)
	client := http.Client{Transport: tp, Timeout: t.requestTimeout}
	res, err := client.Get("http://" + url.Host)
	duration = tp.Duration()
	if err == nil {
		status = res.Status
	} else if duration >= t.requestTimeout {
		status = "TIMEOUT"
	} else {
		status = formatError(err, url)
	}

	return
}

func formatError(err error, url *url.URL) string {
	m := regexp.MustCompile(fmt.Sprintf(`(Get \"http://%s\": )|(net/http: )| \(.*\)|(dial .* %s: )`, url.Host, url.Host))
	return m.ReplaceAllString(err.Error(), "")
}
