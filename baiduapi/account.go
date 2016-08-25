package baiduapi

import (
	"fmt"
	"strconv"
)

const (
	acc_productline = "sms"
	acc_version 	= "service"
	acc_service 	= "AccountService"
)

type Account struct {
	cli			*Client
	username	string
}

func NewAccount(token ,password , username string) *Account {	
	cli := NewClient(token, password, username, acc_productline ,acc_version ,acc_service)
	return &Account{cli: cli, username: username}
}

func (this *Account) GetAccountId() (string, error) {

	body := "{\"accountFields\":[\"userId\"]}"

	data,err := this.cli.Execute("getAccountInfo",body)
	if err != nil {
		fmt.Println("getAccountInfo error! username:",this.username,", error:",err)
		return "",err
	}
	// fmt.Println("data=",data)

	floatUid := data["body"].(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})["userId"].(float64)
	userId := strconv.Itoa(int(floatUid))
	fmt.Println("userId:",userId)
	return userId,nil
}