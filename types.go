package main

type ResponseFailure struct {
	Errcode int    `json:"errcode,omitempty"`
	Errmsg  string `json:"errmsg,omitempty"`
}

func (resp *ResponseFailure) React() {
	if resp.Errcode < 0 {

	}
}
