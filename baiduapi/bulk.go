package baiduapi

import (
	"fmt"	
	"errors"
)

const (
	productline = "sms"
	version 	= "service"
	service 	= "BulkJobService"
)


type Bulk struct {
	cli		*Client
	username	string
}

func NewBulk(token ,password , username string) *Bulk {	
	cli := NewClient(token, password, username, productline ,version ,service)
	return &Bulk{cli: cli, username: username}
}

func (this *Bulk) GetCampaign(accountId string) (string,error){
	body := "{\"accountIds\":["+accountId+"],\"campaignFields\":[\"all\"]}"
	return this.GetAllObjects(body)	
}

//key: accountIds,campaignIds
func (this *Bulk) GetKeyword(ids, key string) (string, error) {
	body := "{\""+key+"\":["+ids+"],\"keywordFields\":[\"all\"],\"includeTemp\":\"false\"}"
	return this.GetAllObjects(body)
}

func (this *Bulk) GetAllObjects(body string) (string,error) {
	// campaign := ","
	// for _, id range campaignIds {
	// 	campaign += string(id)
	// }
	// body := "{\"campaignIds\":["+campaignid+"],\"keywordFields\":[\"all\"],\"includeTemp\":\"false\"}"

	data,err := this.cli.Execute("getAllObjects",body)
	if err != nil {
		fmt.Println("getAllObjects error! username:",this.username,", error:",err)
		return "",err
	}
	// fileId := data["body"].(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})["fileId"].(string)
	// fmt.Println("fileId:",fileId)
	v,err := this.jsonData(data, "fileId")
	if err != nil {
		fmt.Println(err)
		fmt.Println("data=",data)
		return "",err
	}
	fileId := v.(string)
	return fileId,nil
}

func (this *Bulk) GetFileStatus(fileId string) (int,error) {
	body := "{\"fileId\":\""+fileId+"\"}"
	data,err := this.cli.Execute("getFileStatus",body)
	if err != nil {
		fmt.Println("getFileStatus error! username:",this.username,", error:",err)
		return -1,err
	}	
	// status := data["body"].(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})["isGenerated"].(float64)
	v,err := this.jsonData(data, "isGenerated")
	if err != nil {
		fmt.Println(err)
		fmt.Println("data=",data)
		return -1,err
	}
	status := v.(float64)
	// fmt.Println("getFileStatus===",data)
	// fmt.Println("status====",status)
	return int(status),nil
}

//fileKeyname keywords->keywordFilePath, campaign->campaignFilePath
func (this *Bulk) GetFilePath(fileId, fileKeyName string) (string,error) {
	body := "{\"fileId\":\""+fileId+"\"}"
	data,err := this.cli.Execute("getFilePath",body)
	if err != nil {
		fmt.Println("getFilePath error! username:",this.username,", error:",err)
		return "",err
	}
	// fileUrl := data["body"].(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})[fileKeyName].(string)
	v,err := this.jsonData(data, fileKeyName)
	if err != nil {
		fmt.Println(err)
		fmt.Println("data=",data)
		return "",err
	}
	fileUrl := v.(string)
	fmt.Println("fileUrl====",fileUrl)
	return fileUrl,nil
}

func (this *Bulk) jsonData(data map[string]interface{}, keyName string) (interface{}, error) {
	b,ok := data["body"]
	if !ok {
		return nil,errors.New("body is null")
	}
	if b == nil {
		return nil,errors.New("body is nil")
	}

	mb := b.(map[string]interface{})

	d,ok := mb["data"]

	if !ok {
		return nil,errors.New("data is null")
	}

	if d == nil {
		return nil,errors.New("data is nil")
	}

	ad := d.([]interface{})

	if len(ad) < 1 {
		return nil, errors.New("data len is 0")
	}

	md := ad[0].(map[string]interface{})

	v,ok := md[keyName]

	if !ok {
		return nil,errors.New(keyName+" is null")
	}

	if v == nil {
		return nil, errors.New(keyName+" is nil")
	}

	return v,nil
}

