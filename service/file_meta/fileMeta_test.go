package service

import (
	"DisHub/config"
	"fmt"
	"testing"
)

func TestFileMetaService(t *testing.T) {
	// 初始化 FileMetaService
	config.DefaultCfg.MYSQLADDR = "ossfile:ossfile@tcp(localhost:3306)/oss_fileMeta"
	fms := NewFileMetaService()
	defer fms.Close()

	// 插入一个文件元数据
	fm := &FileMetadata{
		Filename: "test.txt",
		Size:     1024,
		Hash:     "abcdef123456",
	}
	inserted := fms.PutMetaData(fm.Filename, fm.Size, fm.Hash)
	if inserted != nil {
		t.Error("Insert failed")
	}

}

func TestGetFileData(t *testing.T) {
	config.DefaultCfg.MYSQLADDR = "ossfile:ossfile@tcp(localhost:3306)/oss_fileMeta"
	fms := NewFileMetaService()
	defer fms.Close()
	f, e := fms.GetMetaData("aaa")
	fmt.Println(f)
	fmt.Println(e)
}
