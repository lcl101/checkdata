package config

const (
	Conf_baidu_url 			= "https://api.baidu.com"
	Conf_baidu_action    	= "API-SDK"
	Conf_baidu_target		= ""
	Conf_baidu_accessToken	= ""

	Tmp_dir		= "tmp"
	Data_dir	= "data"
)

type FileMap struct {
	Key		int
	Index	[]int
}

func NewBaiduFileMap() *FileMap {
	index := make([]int, 3)
	index[0] 	= 0
	index[1] 	= 1
	index[2]	= 3		
	// index[3]	= 9

	return &FileMap{Key: 2, Index: index}
}

func NewLocalFileMap() *FileMap {
	index := make([]int, 3)
	index[0] 	= 3
	index[1] 	= 5
	index[2]	= 0
	// index[3] 	= 5
	

	return &FileMap{Key: 7, Index: index}
}

var (
	BaiduFileMap = NewBaiduFileMap()
	LocalFileMap = NewLocalFileMap()
	IsDebug		= false
)
