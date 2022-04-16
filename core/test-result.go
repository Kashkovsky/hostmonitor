package core

import "net/url"

const (
	StatusOK          = "OK"
	StatusErr         = "Error"
	StatusErrResponse = "ErrorResponse"
)

type TestResult struct {
	Id           string `json:"id"`
	InProgress   bool   `json:"inProgress"`
	url          url.URL
	Tcp          string `json:"tcp"`
	HttpResponse string `json:"httpResponse"`
	Duration     string `json:"duration"`
	Status       string `json:"status"`
}
