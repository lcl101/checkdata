package util

import(
	"github.com/lcl101/checkdata/config"
	"fmt"
	"io"
	"net/http"
	"crypto/tls"
	"os"
	"compress/gzip"
	"bufio"
	"strings"
	"crypto/md5"
	"encoding/hex"
)

func IsDebug() bool {
	return true
}

func HandleError(err error) {
	fmt.Println("occurred error:", err)    
    panic("error.........")
}

func MD5(msg string) string {
	if msg == "" {
		fmt.Println("run MD5, msg is null")
		return ""
	}	
	md5Ctx := md5.New()
    md5Ctx.Write([]byte(msg))
    cipherStr := md5Ctx.Sum(nil)
    str := hex.EncodeToString(cipherStr)
    // fmt.Println(str)
    return str
}

func DownloadFile(url, fileName string) error {
	tr := &http.Transport{
        TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
        DisableCompression: true,
    }
    client := &http.Client{Transport: tr}
	res, err := client.Get(url)  
    if err != nil {  
        return err
    }
    defer res.Body.Close()
    f, err := os.Create(fileName)  
    if err != nil {  
        return err
    }
    defer f.Close() 
    io.Copy(f, res.Body) 

    return nil
}

func UnpackGzipFile(gzFilePath, dstFilePath string) (int64, error) {
    gzFile, err := os.Open(gzFilePath)
    if err != nil {
        return 0, fmt.Errorf("Failed to open file %s for unpack: %s", gzFilePath, err)
    }
    dstFile, err := os.OpenFile(dstFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
    if err != nil {
        return 0, fmt.Errorf("Failed to create destination file %s for unpack: %s", dstFilePath, err)
    }

    ioReader, ioWriter := io.Pipe()

    go func() { // goroutine leak is possible here
        gzReader, _ := gzip.NewReader(gzFile)
        // it is important to close the writer or reading from the other end of the
        // pipe or io.copy() will never finish
        defer func(){
            gzFile.Close()
            gzReader.Close()
            ioWriter.Close()
        }()

        io.Copy(ioWriter, gzReader)
    }()

    written, err := io.Copy(dstFile, ioReader)
    if err != nil {
        return 0, err // goroutine leak is possible here
    }
    ioReader.Close()
    dstFile.Close()

    return written, nil
}

func LocalFile2DB(fileName string, db *Leveldb) error {
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
        arr := strings.Split(line,",")
        kwid := arr[config.LocalFileMap.Key]
        value := ""
        for _, r := range config.LocalFileMap.Index {
        	value += ","+arr[r]
        }
        fmt.Println("kwid=",kwid," || value=",value)
        db.Put(kwid,MD5(value))
    }
    return nil
}

