package baiduapi

import (
	"github.com/lcl101/checkdata/util"
	"encoding/json"
)

type Client struct {
	conf	*util.Config
	con 	*util.Conn
}

func NewClient(token ,password , username , productline ,version ,service string) *Client {
	conf := util.NewConfig(token ,password , username , productline ,version ,service)
	con := util.NewConn(conf)

	return &Client{conf: conf, con: con}
}

func (this *Client) Execute(method string,body string) (map[string]interface{},error) {
	jsonStr, err := this.con.Do(method,body)
	if err != nil {
		// fmt.Println(err)
		return nil,err
	}
	var data map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		// fmt.Println(err)
		return nil,err
	}
	return data,err
}