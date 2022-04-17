package core

import (
	"io/ioutil"
	"strings"
)

type WatchConfig struct {
	ConfigUrl      string
	TestInterval   int
	RequestTimeout int
	UpdateInterval int
}

func (c *WatchConfig) UpdateURLs() (string, error) {
	if strings.Contains(c.ConfigUrl, "://") {
		return GetStringFromURL(c.ConfigUrl)
	}

	raw, err := ioutil.ReadFile(c.ConfigUrl)

	if err != nil {
		return "", err
	}

	return string(raw), nil
}
