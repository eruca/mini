package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	create_menu_url = "https://api.weixin.qq.com/cgi-bin/menu/create"
)

//go:embed menu.json
var menu_json string

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
