package util

import (
	"fmt"
	"time"

	"github.com/imroc/req/v3"

	"weavatar/pkg/wangsu/common/model"
)

func Call(requestMsg model.HttpRequestMsg) string {
	client := req.C()
	client.SetTimeout(10 * time.Second)
	client.SetCommonRetryCount(2)
	client.ImpersonateSafari()
	client.EnableDumpAll()
	fmt.Println(requestMsg)

	request := client.R()
	request.SetHeaders(requestMsg.Headers)

	resp, _ := client.R().SetHeaders(requestMsg.Headers).SetQueryParams(requestMsg.Params).SetBody(requestMsg.Body).Send(requestMsg.Method, requestMsg.Url)

	return resp.String()
}
