package model

type HttpRequestMsg struct {
	Uri           string
	Url           string
	Method        string
	Host          string
	Params        map[string]strings
	Headers       map[string]string
	Body          string
	SignedHeaders string
	Msg           string
}
