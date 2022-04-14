package core

type WatchConfig struct {
	ConfigUrl      string
	TestInterval   int
	RequestTimeout int
	UpdateInterval int
}

func (c *WatchConfig) UpdateURLs() (string, error) {
	return GetStringFromURL(c.ConfigUrl)
}
