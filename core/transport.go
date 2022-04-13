package core

import (
	"net"
	"net/http"
	"time"
)

type Transport struct {
	rtp       http.RoundTripper
	dialer    *net.Dialer
	connStart time.Time
	connEnd   time.Time
	reqStart  time.Time
	reqEnd    time.Time
}

func NewTransport(timeout time.Duration) *Transport {
	tr := &Transport{
		dialer: &net.Dialer{
			Timeout:   timeout,
			KeepAlive: timeout,
		},
	}
	tr.rtp = &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		Dial:                tr.Dial,
		TLSHandshakeTimeout: timeout,
	}
	return tr
}

func (tr *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	tr.reqStart = time.Now()
	resp, err := tr.rtp.RoundTrip(r)
	tr.reqEnd = time.Now()
	return resp, err
}

func (tr *Transport) Dial(network, addr string) (net.Conn, error) {
	tr.connStart = time.Now()
	cn, err := tr.dialer.Dial(network, addr)
	tr.connEnd = time.Now()
	return cn, err
}

func (tr *Transport) ReqDuration() time.Duration {
	return tr.Duration() - tr.ConnDuration()
}

func (tr *Transport) ConnDuration() time.Duration {
	return tr.connEnd.Sub(tr.connStart)
}

func (tr *Transport) Duration() time.Duration {
	return tr.reqEnd.Sub(tr.reqStart)
}
