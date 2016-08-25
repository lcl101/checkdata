package main 

import (
	"github.com/lcl101/checkdata/baiduapi"
	"github.com/lcl101/checkdata/util"
	"github.com/lcl101/checkdata/config"
	"github.com/c4pt0r/ini"
	"strconv"
	"fmt"
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

func test_downcsv(token, password, username string) {
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

	baiduapi.DownCsv(bulk , id, key, gzFileName, csvFileName)
}

func main() {
	conf := ini.NewConf("./conf.ini")	
	count := conf.Int("count","count",-1)
	conf.Parse()
	if *count < 1 {
		fmt.Println("conf count is ", *count)
	}

	for i:=1; i<= *count; i++ {
		str := strconv.Itoa(i)
		token := conf.String("a_"+str,"token", "")
		password := conf.String("a_"+str,"password", "")
		username := conf.String("a_"+str,"username", "")
		conf.Parse()
		fmt.Println("==================begin down csv: "+*username)
		test_downcsv(*token, *password, *username)
		fmt.Println("==================end down csv: "+*username)		
	}
}