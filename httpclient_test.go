package httpclient

import (
	"fmt"
	"testing"
)

func Test_getUrl(t *testing.T) {
	var (
		u string =  "https://oapi.dingtalk.com/robot/send"
		paramkey string = "access_token"
		paramvalue string = "ac0066b0a6335d82953934adf6dc32e4d8f539851b7bf788f14e0d744a35e5a1"
		queryparam map[string]string = make(map[string]string,0)
	)
	
	queryparam[paramkey] = paramvalue
	client := NewHttpClient(10)
	
	url,err := client.SetQuery(u,queryparam)
	if err != nil {
		fmt.Printf("url set query param failed:%v\n",err)
		return
	}
	
	fmt.Printf("url:%v\n",url)
}
