package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	create_menu_url = "https://api.weixin.qq.com/cgi-bin/menu/create"
	menu_json       = `{
		"button":[
		{	
			 "type":"click",
			 "name":"今日歌曲",
			 "key":"V1001_TODAY_MUSIC"
		 },
		 {
			  "name":"菜单",
			  "sub_button":[
			  {	
				  "type":"view",
				  "name":"搜索",
				  "url":"http://www.soso.com/"
			   },
			   {
				  "type":"click",
				  "name":"赞一下我们",
				  "key":"V1001_GOOD"
			   }]
		  }]
	}`
)

func CreateMenu() error {
	client := http.Client{}

	access_token, err := GetTokenAccess()
	if err != nil {
		return err
	}

	reader := strings.NewReader(menu_json)

	resp, err := client.Post(fmt.Sprintf("%s?access_token=%s", create_menu_url, access_token), "application/json;charset=UTF-8", reader)
	if err != nil {
		return err
	}

	respFail := ResponseFailure{Errcode: -1}
	err = json.NewDecoder(resp.Body).Decode(&respFail)
	if err != nil {
		return err
	}
	log.Printf("respFail: %v", respFail)

	if respFail.Errcode != 0 {
		return errors.New(respFail.Errmsg)
	}
	return nil
}
