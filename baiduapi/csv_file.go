package baiduapi

import (
	"github.com/lcl101/checkdata/util"	
	"github.com/lcl101/checkdata/config"	
	"fmt"
	"time"
	"io"
	"bufio"
	"os"
	"strings"
)

func DownCsv(bulk *Bulk, id, key, gzFileName, csvFileName string) error {
	fileId,err := bulk.GetKeyword(id, key)
	if err != nil {
		return err
	}
	status := 0
	for status < 3 {
		status, err = bulk.GetFileStatus(fileId)
		if err != nil {
			return err
		}
		// fmt.Println("file status is: ",status)
		time.Sleep(time.Second * 3)
	}

	fileUrl, err := bulk.GetFilePath(fileId,"keywordFilePath")
	if err != nil {
		return err
	}
	err = util.DownloadFile(fileUrl,gzFileName)
	if err != nil {
		return err
	}
	count, err := util.UnpackGzipFile(gzFileName,csvFileName)
	if err != nil {
		return err
	}
	fmt.Println("ungzip count: ", count)
	return nil
}

func DownCampCsv(bulk *Bulk, accountId, gzFileName, csvFileName string) error {
	fileId,err := bulk.GetCampaign(accountId)
	if err != nil {
		return err
	}
	status := 0
	for status < 3 {
		status, err = bulk.GetFileStatus(fileId)
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 2)
	}

	fileUrl, err := bulk.GetFilePath(fileId,"campaignFilePath")
	if err != nil {
		return err
	}
	err = util.DownloadFile(fileUrl,gzFileName)
	if err != nil {
		return err
	}
	count, err := util.UnpackGzipFile(gzFileName,csvFileName)
	if err != nil {
		return err
	}
	fmt.Println("ungzip count: ", count)
	return nil
}

func GetCampIds(bulk *Bulk, accountId string) ([]string,error) {
	gzFileName := "./"+config.Tmp_dir+"/acc_"+accountId+".gz"
	csvFileName := "./"+config.Tmp_dir+"/acc_"+accountId+".csv"
	err := DownCampCsv(bulk, accountId, gzFileName,csvFileName)
	if err != nil {
		return nil,err
	}

	f, err := os.Open(csvFileName)  
    defer f.Close() 
    if nil != err {
    	return nil, err
    }

    data := make([]string,0)
     
    buff := bufio.NewReader(f)
    line, err := buff.ReadString('\n') //去掉第一行
    for {  
        line, err = buff.ReadString('\n')  
        if err != nil || io.EOF == err{  
            break  
        }
        // fmt.Println("line====",line)  
        arr := strings.Split(line,"\t")
        data = append(data,arr[0])
    }
    return data,nil
}

func IntoDb(bulk *Bulk, db *util.Leveldb, id, key string) error {
	gzFileName := "./"+config.Tmp_dir+"/"+key+"_"+id+".gz"
	csvFileName := "./"+config.Tmp_dir+"/"+key+"_"+id+".csv"
	err := DownCsv(bulk, id, key, gzFileName, csvFileName)
	if err != nil {
		return err
	}
	return File2DB(csvFileName, db)
}

func File2DB(fileName string, db *util.Leveldb) error {
	f, err := os.Open(fileName)  
    defer f.Close() 
    if nil != err {
    	return err
    }
     
    buff := bufio.NewReader(f)  
    for {  
        line, err := buff.ReadString('\n')  
        if err != nil || io.EOF == err{  
            break  
        }
        // fmt.Println("line====",line)  
        arr := strings.Split(line,"\t")
        kwid := arr[config.BaiduFileMap.Key]
        value := ""
        for _, r := range config.BaiduFileMap.Index {
        	value += ","+arr[r]
        }
        fmt.Println("kwid=",kwid," || value=",value)
        //存入md5
        db.Put(kwid,util.MD5(value))
    }
    return nil
}

