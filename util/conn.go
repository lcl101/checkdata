package util

import (
	"github.com/lcl101/checkdata/config"
	"net/http"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"strings"
	
)

type Config struct {
	url			string
	action		string
	target		string
	accessToken	string

	token		string
	password	string
	username 	string

	productline string
	version		string
	service 	string
}

func (this *Config) toJson() string {
	return " {"	+
	this.add("username",this.username)+"," 	+
	this.add("token",this.token)+","+
	this.add("target",this.target)+"," +
	this.add("accessToken",this.accessToken)+","+
	this.add("action",this.action)+"," +
	this.add("password",this.password) +"}"
}

func (this *Config) add(key, value string) string {
	return "\""+key+"\":\"" + value + "\""
}

func NewConfig(token ,password , username , productline ,version ,service string) *Config {
	return &Config{
		url:			config.Conf_baidu_url,
		action:			config.Conf_baidu_action,
		target:			config.Conf_baidu_target,
		accessToken:	config.Conf_baidu_accessToken,

		token: 			token,
		username: 		username,
		password:		password,
		productline:	productline,
		version:		version,
		service:		service }
}

type Conn struct {
	conf 	*Config
}

func NewConn(conf *Config) *Conn {
	return &Conn{conf: conf}
}

// 使用Do方法，不使用Post方法
func (this *Conn) Post(method string,body string) (string, error) {
	url := this.conf.url + "/json/"+this.conf.productline+"/"+this.conf.version+"/" + this.conf.service + "/" +method
	json := "{ \"body\": "+body+", \"header\":" + this.conf.toJson()+"}"
	fmt.Println("json========",json)
	tr := &http.Transport{
        TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
        DisableCompression: true,
    }

    client := &http.Client{Transport: tr}
	resp, err := client.Post(url,"application/json;charset=utf-8", strings.NewReader(json))
	if err != nil {
		// HandleError(err)
		return "",err
	}
	defer resp.Body.Close()

	fmt.Println("resp=", resp)
	fmt.Println("resp.header=",resp.Header)
	fmt.Println("resp.body=",resp.Body)

	tmpBody, err := ioutil.ReadAll(resp.Body)
	fmt.Println("rrrr;error=",err)
	fmt.Println("body=",string(tmpBody))
	if IsDebug() {
		fmt.Println("resp=", resp)
		fmt.Println("resp.header=",resp.Header)
		fmt.Println("resp.body=",resp.Body)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "",err
	}
	if IsDebug() {	
		fmt.Println("body=",string(data))
	}
	return string(data),nil
	
}

// { "body": {"campaignIds":[66601576],"keywordFields":["all"],"includeTemp":"false"}, "header": 
//{"username":"baidu-京东POP放量三二8161355","token":"42ec02df5e4822e054f2e2e011a41f20","target":"","accessToken":"","action
// ":"API-SDK","password":"JD@pc0719"}}

func (this *Conn) Do(method string,body string) (string, error) {
	url := this.conf.url + "/json/"+this.conf.productline+"/"+this.conf.version+"/" + this.conf.service + "/" +method
	json := "{ \"body\": "+body+", \"header\":" + this.conf.toJson()+"}"
	if IsDebug() {
		fmt.Println("json=",json)
	}
	tmp := ioutil.NopCloser(strings.NewReader(json))
	req, err := http.NewRequest("POST", url, tmp)
	if err != nil {
		// HandleError(err)
		return "",err
	}
	req.Header.Set("Content-type","application/json;charset=utf-8")
	req.Header.Set("Connection","keep-alive")

	tr := &http.Transport{
        TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
        DisableCompression: true,
    }
	// if IsDebug() {
	// 	fmt.Println("req=",req)
	// }
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		// HandleError(err)
		return "",err
	}
	defer resp.Body.Close()
	if IsDebug() {
		fmt.Println("resp=", resp)
		fmt.Println("resp.header=",resp.Header)
		fmt.Println("resp.body=",resp.Body)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "",err
	}
	if IsDebug() {	
		fmt.Println("body=",string(data))
	}
	return string(data),nil
}
