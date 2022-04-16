package web

import "github.com/Kashkovsky/hostmonitor/core"

const (
	KindResult = "result"
	KindReset  = "reset"
)

type Message struct {
	Kind string      `json:"kind"`
	Data interface{} `json:"data"`
}

func NewResultMessage(res core.TestResult) Message {
	return Message{
		Kind: KindResult,
		Data: res,
	}
}

func NewResetMessage() Message {
	return Message{
		Kind: KindReset,
	}
}
