package main 

import (
	"github.com/lcl101/checkdata/baiduapi"
	"github.com/lcl101/checkdata/util"
	"github.com/lcl101/checkdata/config"
	"github.com/c4pt0r/ini"
	"github.com/tealeg/xlsx"
	"strconv"
	"fmt"
	"os"
	"bufio"
	"io"
	"flag"
)

func test(){
	token := ""
	password := ""
	username := "baidu-京东POP放量三二8161355"
	//以上信息需要放入配置文件

	acc := baiduapi.NewAccount(token ,password , username)	
	accountId, err := acc.GetAccountId()
	if err != nil {
		util.HandleError(err)
	}	

	localFileName := "./local/"+accountId+".csv"

	bulk := baiduapi.NewBulk(token ,password , username)
	//目前以account维度去check关键字
	id := accountId
	key := "accountIds"
	// key := "campaignIds"

	//生成数据库
	localdb := util.NewDB("local_"+id)
	defer localdb.Close()
	accdb := util.NewDB("acc_"+id)
	defer accdb.Close()
	//下载csv文件并入库
	err = baiduapi.IntoDb(bulk, accdb, id, key)
	if err != nil {
		util.HandleError(err)
	}
	err = util.LocalFile2DB(localFileName, localdb)
	if err != nil {
		util.HandleError(err)
	}
}

func test_downcsv(token, password, username string) (string,string) {
	//以上信息需要放入配置文件
	acc := baiduapi.NewAccount(token ,password , username)	
	accountId, err := acc.GetAccountId()
	if err != nil {
		util.HandleError(err)
	}	

	bulk := baiduapi.NewBulk(token ,password , username)
	//目前以account维度去check关键字
	id := accountId
	key := "accountIds"
	// key := "campaignIds"

	gzFileName := "./"+config.Tmp_dir+"/"+key+"_"+id+ "_"+username+".gz"
	csvFileName := "./"+config.Tmp_dir+"/"+key+"_"+id+ "_"+username+".csv"

	err = baiduapi.DownCsv(bulk , id, key, gzFileName, csvFileName)
	if err != nil {
		fmt.Println("username=",username,", error=",err)
		util.HandleError(err)
	}
	return csvFileName,accountId
}

func merge(subCsv string, w *bufio.Writer) (int,error) {
	f, err := os.Open(subCsv)  
    defer f.Close() 
    if nil != err {
    	return -1,err
    }
     
    buff := bufio.NewReader(f)  
    //除去第一行
    line, err := buff.ReadString('\n')
    if nil != err {
    	return -1,err
    }
    i := 0
    for {  
        line, err = buff.ReadString('\n') 
        if err != nil || io.EOF == err{  
            break 
        }
        if nil != err {
        	return -1,err
        }        
        _,err = w.WriteString(line)
        if nil != err {
        	return -1,err 
        }
        i++
    }
    err = w.Flush()
    if nil != err {
    	return -1,err
    }
    return i, nil
}

func file_merge() {
	conf := ini.NewConf("./conf.ini")	
	count := conf.Int("count","count",-1)


	conf.Parse()
	if *count < 1 {
		fmt.Println("conf count is ", *count)
	}
	csv := "./"+config.Tmp_dir+"/data.csv"
	csvf, err := os.OpenFile(csv, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
    defer csvf.Close() 
    if nil != err {
    	util.HandleError(err)
    }
    w := bufio.NewWriter(csvf)

    total := 0
	for i:=1; i<= *count; i++ {
		str := strconv.Itoa(i)
		token := conf.String("a_"+str,"token", "")
		password := conf.String("a_"+str,"password", "")
		username := conf.String("a_"+str,"username", "")
		conf.Parse()
		fmt.Println("==================begin csv: "+*username)
		subCsv,_ := test_downcsv(*token, *password, *username)
		fmt.Println("==================end down csv: "+*username)
		c, err := merge(subCsv, w)
		if nil != err {
			util.HandleError(err)
		}
		fmt.Println("username=",*username,",lines=",c)
		total = total + c
		fmt.Println("==================end merge csv: "+*username)
	}
	fmt.Println("total lines=",total)
}

func count(token, password, username string, index int, c1 chan int) (int,error){
	fmt.Println("downloading ======",username,", index=",index)
	csv,accountId := test_downcsv(token, password, username)
	fmt.Println("opening ======",username,", index=",index)

	f, err := os.Open(csv)  
    defer f.Close() 
    if nil != err {
    	return -1,err
    }
    
    buff := bufio.NewReader(f)    
    i := -1
    for {  
        _, err := buff.ReadString('\n') 
        if err != nil || io.EOF == err{  
            break 
        }         
        i++
    }       
    fmt.Println("opened ======",username,", index=",index)
    fmt.Println(accountId,",",i,",",username,",MMMARKMM",", index=",index) 
    c1 <- 1
    return i,nil
}

func file_count() {
	conf := ini.NewConf("./jd.ini")	
	xls := conf.String("jd","xls","./jd.xlsx")
	tIndex := conf.Int("xls","tIndex",-1)
	pIndex := conf.Int("xls","pIndex",-1)
	uIndex := conf.Int("xls","uIndex",-1)
	start := conf.Int("xls","start",-1)
	end   := conf.Int("xls","end",-1)
	conf.Parse()
	if *tIndex <0 || *pIndex <0 || *uIndex<0 {
		fmt.Println("index error!")
		return
	}
	if *start<0 || *end<0 {
		fmt.Println("start,end error")
		return
	}
	xlFile, err := xlsx.OpenFile(*xls)
    if err != nil {
        util.HandleError(err)
    }   
    c1 := make(chan int) 
    threadCount := 0
    sheet := xlFile.Sheets[0]
    for i := *start; i<*end; i++ {
    	row := sheet.Rows[i]
    	token := row.Cells[*tIndex].String()
    	password := row.Cells[*pIndex].String()
    	username := row.Cells[*uIndex].String()
    	fmt.Println("token=",token,",pd=",password,"un=",username)    	
    	go count(token, password, username, i, c1)
    	threadCount++
    	if threadCount>9 {
	    	for threadCount>0{
	    		fmt.Println("--------------",threadCount)	    		
	    		<- c1
	    		threadCount--
	    	}
	    }
    }
    for threadCount>0{
		fmt.Println("--------------",threadCount)		
		<- c1
		threadCount--
	}
}

func main() {
	method := flag.Int("cmd", 0, "cmd")
	debug := flag.Int("d",0,"debug")
	flag.Parse()
	fmt.Println("cmd=",*method)
	if *debug == 1 {
		config.IsDebug = true
	}
	if 1 == *method {
		// file_merge()
	} else if 2 == *method {
		file_count()
	} else {
		fmt.Println("method error. method=",*method)
	}	
}