package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/karlseguin/ccache/v2"
)

// token access url
const (
	token_access_url = "https://api.weixin.qq.com/cgi-bin/token"
	appid            = "wx4eed36199e7795df"
	secret           = "30022d0725786808833d8a7fbe3238e7"

	token_access = "token_access"

	errcode1 = `{"errcode":`
	errcode2 = `{"errmsg":`
)

var (
	cache          = ccache.New(ccache.Configure())
	GetTokenAccess = getTokenAccessWrapper(appid, secret)
)

func getTokenAccessWrapper(appid, secret string) func() (string, error) {
	expiresIn := 7200 * time.Second

	return func() (string, error) {
		item, err := cache.Fetch(token_access, expiresIn, func() (interface{}, error) {
			req := TokenAccessRequest{Appid: appid, Secret: secret}
			ta, err := req.Do()
			if err != nil {
				return nil, err
			}
			expiresIn = time.Duration(ta.ExpiresIn) * time.Second
			return ta.AccessToken, nil
		})
		if err != nil {
			return "", err
		}

		return item.Value().(string), nil
	}
}

type TokenAccessRequest struct {
	Appid  string `json:"appid,omitempty"`
	Secret string `json:"secret,omitempty"`
}

func (t *TokenAccessRequest) url() string {
	return fmt.Sprintf(
		"%s?grant_type=client_credential&appid=%s&secret=%s",
		token_access_url, t.Appid, t.Secret)
}

func (t *TokenAccessRequest) Do() (*TokenAccess, error) {
	errCode := -1

	for i := 0; errCode == -1; i++ {
		resp, err := http.Get(t.url())
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// 如果以errcode1或errcode2开始就是返回错误
		if bytes.HasPrefix(data, []byte(errcode1)) || bytes.HasPrefix(data, []byte(errcode2)) {
			var respFail ResponseFailure
			err = json.Unmarshal(data, &respFail)
			if err != nil {
				return nil, err
			}
			if respFail.Errcode == -1 {
				errCode = -1
				time.Sleep(time.Duration(i+1) * time.Microsecond * 500)
				continue
			}
			return nil, errors.New(respFail.Errmsg)
		}

		ta := &TokenAccess{}
		err = json.Unmarshal(data, ta)
		if err != nil {
			return nil, err
		}
		return ta, nil
	}

	return nil, nil
}

type TokenAccess struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
}
