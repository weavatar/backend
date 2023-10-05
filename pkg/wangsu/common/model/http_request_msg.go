package model

type HttpRequestMsg struct {
	Uri           string
	Url           string
	Method        string
	Host          string
	Params        string
	Headers       map[string]string
	Body          string
	SignedHeaders string
	Msg           string
}
