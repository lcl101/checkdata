package util

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/lcl101/checkdata/config"
	"fmt"
)

type Leveldb struct {
	dbFile		string
	db 			*leveldb.DB
}

func NewDB(name string) *Leveldb {
	dbFile := "./"+config.Data_dir+"/"+name+".db"

	db, err := leveldb.OpenFile(dbFile, nil)

	if err != nil {
		fmt.Println("db open file error: ", err)
		return nil
	}

	return &Leveldb{dbFile: dbFile, db: db}
}

func (this *Leveldb) Close() {
	if this.db != nil {
		this.db.Close()
	}
}

func (this *Leveldb) Put(kwId, value string) {
	err := this.db.Put([]byte(kwId), []byte(value), nil)
	if err != nil {
		fmt.Println("pub error: ", err)
	}
}

func (this *Leveldb) Del(kwId string) {
	err := this.db.Delete([]byte(kwId), nil)
	if err != nil {
		fmt.Println("del error: ", err)
	}
}

func (this *Leveldb) Get(kwId string) string {
	data, err := this.db.Get([]byte(kwId), nil)
	if err != nil {
		fmt.Println("del error: ", err)
		return ""
	}
	return string(data)
}