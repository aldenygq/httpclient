package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultDialTimeout             = 2 * time.Second
	defaultKeepAlive               = 2 * time.Second
)

var HttpClient *http.Client


func NewHttpClient() *http.Client {
	HttpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   defaultDialTimeout,
				KeepAlive: defaultKeepAlive,
			}).DialContext,
		},
		Timeout: defaultDialTimeout,
	}
	return HttpClient
}

func DoReq(method, u string, body interface{},header map[string]string,queryparams map[string]string) (*http.Response, error) {
	var (
		req *http.Request
		err error
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()
	
	req,err = NewReqByMethod(ctx,method, u,body,queryparams)
	if err != nil {
		return nil, err
	}
	
	if header != nil && len(header) > 0{
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	req = req.WithContext(ctx)
	
	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil,err
	}
	
	//defer resp.Body.Close()
	
	return resp, nil
}
func NewReqByMethod(ctx context.Context,method, u string, body interface{},queryparams map[string]string) (*http.Request,error) {
	var (
		err error
		req *http.Request = &http.Request{}
		b []byte = make([]byte,0)
	)
	if body != nil {
		b, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	if method == "GET" || method == "DELETE" {
		req, err = http.NewRequestWithContext(ctx,method, u, nil)
		_,err = SetQuery(u,queryparams)
		if err != nil {
			return nil,err
		}
	} else if  method == "POST" || method == "PUT" {
		req, err = http.NewRequestWithContext(ctx,method, u, bytes.NewReader(b))
	} else {
		return nil,errors.New("request method invalid")
	}
	return req,nil
}
//set query
func SetQuery(u string, params map[string]string) (string, error) {
	var q url.Values
	_u, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	if params != nil {
		q = _u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
	}
	_u.RawQuery = q.Encode()
	
	return _u.String(), nil
}

//去掉url配置中最后一个‘/’
func CutLastestSlash(u string) string {
	if u[len(u)-1] == '/' {
		u = u[0:(len(u) - 1)]
	}
	
	return u
}