package main

type ResponseFailure struct {
	Errcode int    `json:"errcode,omitempty"`
	Errmsg  string `json:"errmsg,omitempty"`
}

type ScanCodeInfo struct {
	ScanType   string `xml:"ScanType"`
	ScanResult string `xml:"ScanResult"`
}

type WxMenuEvent struct {
	ToUserName   string        `xml:"ToUserName"`
	FromUserName string        `xml:"FromUserName"`
	CreateTime   int64         `xml:"CreateTime"`
	MsgType      string        `xml:"MsgType"`
	Event        string        `xml:"Event"` //VIEW
	EventKey     string        `xml:"EventKey"`
	MenuId       string        `xml:"MenuId"`
	ScanCodeInfo *ScanCodeInfo `xml:"ScanCodeInfo"` //专属于扫码
	Content      string        `xml:"Content"`
	Ticket       string        `xml:"Ticket"`
}

func (e *WxMenuEvent) Dispatch() {
	switch e.MsgType {
	case "VIEW":
	case "CLICK":
	case "SUBSCRIBE":
	case "SCAN":

	}
}
