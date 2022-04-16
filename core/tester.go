package core

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

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
			tcp := "-"
			t.out <- TestResult{
				Id:         url.String(),
				InProgress: true,
				Tcp:        tcp,
			}

			if url.Scheme == "tcp" {
				pass := t.tcp(url)
				tcp = fmt.Sprintf("%d/10", pass)
			}
			response, duration, status := t.http(url)
			t.out <- TestResult{
				Id:           url.String(),
				url:          *url,
				Tcp:          tcp,
				HttpResponse: response,
				Duration:     strconv.FormatInt(duration.Milliseconds(), 10) + "ms",
				Status:       status,
			}
			time.Sleep(t.testInterval)
		}
	}
}

func (t *Tester) tcp(url *url.URL) int {
	tp := NewTransport(time.Second)
	pass := 0
	for i := 0; i < 10; i++ {
		_, err := tp.Dial(url.Scheme, url.Host)
		if err == nil {
			pass++
		}
	}

	return pass
}

func (t *Tester) http(url *url.URL) (statusMessage string, duration time.Duration, status string) {
	tp := NewTransport(t.requestTimeout)
	client := http.Client{Transport: tp, Timeout: t.requestTimeout}
	addr := strings.Replace(url.String(), "tcp", "http", 1)
	res, err := client.Get(addr)
	duration = tp.Duration()
	if err == nil {
		statusMessage = res.Status
		if res.StatusCode >= 500 {
			status = StatusErrResponse
		} else {
			status = StatusOK
		}
	} else if duration >= t.requestTimeout {
		statusMessage = "TIMEOUT"
		status = StatusErrResponse
	} else {
		statusMessage = formatError(err, url)
		status = StatusErr
	}

	return
}

func formatError(err error, url *url.URL) string {
	parts := strings.Split(err.Error(), ":")
	return parts[len(parts)-1]
}
