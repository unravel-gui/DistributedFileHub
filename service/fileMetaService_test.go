package service

import (
	"DisHub/common"
	"DisHub/common/db"
	"DisHub/common/utils"
	"DisHub/repository"
	"github.com/google/uuid"
	"testing"
)

func Test_FileMeta(t *testing.T) {
	clearTestData(t)
	fms := NewFileMetaService(common.TEST_MYSQL_ADDR)
	uuids := uuid.New()
	uid := common.TEST_UID
	str := uuids.String()
	fm := &repository.FileMeta{
		Hash:        str,
		Uid:         uid,
		Dir:         0,
		Name:        str,
		ContentType: str,
		UploadTime:  utils.GetNow(),
	}
	_, err := fms.PutFileMeta(fm)
	if err != nil {
		t.Errorf("put file meta err:%v", err)
	}
	fm1, err := fms.GetFileMetaByHash(uid, str)
	if err != nil {
		t.Errorf("get file meta err:%v", err)
	}
	if fm1 == nil {
		t.Errorf("get empty file meta")
	}
	if fm1.Uid != fm.Uid || fm1.Hash != fm.Hash {
		t.Errorf("get diff file meta, orgin:%+v, get:%+v", fm, fm1)
	}

	ok, err := fms.DeleteFileMetaByHash(uid, str)
	if err != nil {
		t.Errorf("del file meta err:%v", err)
	}
	if !ok {
		t.Errorf("del failed")
	}
	ok, err = fms.DeleteFileMetaByHash(uid, str)
	if err != nil {
		t.Errorf("del file meta err:%v", err)
	}
	if ok {
		t.Errorf("duplice del should be failed")
	}
	fm1, err = fms.GetFileMetaByHash(uid, str)
	if err != nil {
		t.Errorf("get file meta err:%v", err)
	}
	if fm1 != nil {
		t.Errorf("get empty file meta")
	}
}

func Test_FileMetaDir(t *testing.T) {
	clearTestData(t)
	fms := NewFileMetaService(common.TEST_MYSQL_ADDR)
	uuid := uuid.New()
	str := uuid.String()
	uid := common.TEST_UID
	fm := &repository.FileMeta{
		Hash:        str,
		Uid:         uid,
		Dir:         0,
		Name:        str,
		ContentType: common.FOLDER,
		UploadTime:  utils.GetNow(),
		UpdateTime:  utils.GetNow(),
	}
	// /中新建目录
	fms.PutFileMeta(fm)
	// 第一个目录
	fid := fm.Fid
	fm.Fid = 0
	// /中新建目录
	fm.Name = uuid.String()
	fm.Fid = 0
	fms.PutFileMeta(fm)
	// /中新建文件
	fm.ContentType = uuid.String()
	fm.Fid = 0
	fms.PutFileMeta(fm)
	// /中新建文件
	fm.Name = uuid.String()
	fm.Fid = 0
	fms.PutFileMeta(fm)

	// 二级目录中创建两个目录两个文件
	fm.Dir = fid
	fm.Name = "4"
	fm.ContentType = common.FOLDER
	fm.Fid = 0
	fms.PutFileMeta(fm)
	// 第二层目录fid
	fid2 := fm.Fid
	// 新建文件
	fm.Name = "2"
	fm.ContentType = uuid.String()
	fm.Fid = 0
	fms.PutFileMeta(fm)
	fm.Name = "1"
	fm.Fid = 0
	fms.PutFileMeta(fm)
	// 中新建文件
	fm.Name = "3"
	fm.ContentType = common.FOLDER
	fm.Fid = 0
	fms.PutFileMeta(fm)

	// 三级目录，创建一个目录一个文件
	fm.Dir = fid2
	fm.Fid = 0
	fms.PutFileMeta(fm)
	fm.ContentType = uuid.String()
	fm.Fid = 0
	fms.PutFileMeta(fm)

	// 查询二级目录的内容
	fileMetas, err := fms.GetFileMetasByUserAndDir(uid, fid)
	if err != nil {
		t.Fatalf("get file meta by dir failed, err:%v", err)
	}
	files := []string{"3", "4", "1", "2"}
	for i, file := range files {
		if fileMetas[i].Name != file {
			t.Fatalf("get wrong file meta, %+v", fileMetas)
		}
	}
	err = fms.LogicalDeleteFileMetasByDir(uid, []int{fid})
	if err != nil {
		t.Fatalf("del file meta by dir failed, err:%v", err)
	}
	fileMetas, err = fms.GetFileMetasByUserAndDir(uid, fid)
	if err != nil {
		t.Fatalf("get file meta by dir failed, err:%v", err)
	}
	if len(fileMetas) != 0 {
		t.Fatalf("del file meta by dir failed, len(fileMetas)=:%d", len(fileMetas))
	}
}

func clearTestData(t *testing.T) {
	db, err := db.NewConnect(common.TEST_MYSQL_ADDR)
	if err != nil {
		t.Fatalf("clear test data failed")
	}
	result := db.Where("uid = ?", common.TEST_UID).Delete(&repository.FileMeta{})
	if result.Error != nil {
		t.Fatalf("clear test data failed,err:%v\n", err)
	}
}
