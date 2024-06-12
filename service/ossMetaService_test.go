package service

import (
	"DisHub/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOSSMetaService(t *testing.T) {
	// 初始化测试用的 OSSMetaService
	oms := NewOSSMetaService(common.TEST_MYSQL_ADDR)
	defer oms.Close()
	hash := uuid.New().String()
	// 测试 PutMetaData 方法
	err := oms.PutMetaData(hash, 1024)
	assert.NoError(t, err, "PutMetaData should not return error")
	err = oms.PutMetaData(hash, 1024)
	assert.Error(t, err, "PutMetaData should return error")

	// 测试 GetMetaData 方法
	meta, err := oms.GetMetaData(hash)
	assert.NoError(t, err, "GetMetaData should not return error")
	assert.NotNil(t, meta, "GetMetaData should return metadata")

	// 测试 DelMetaData 方法
	err = oms.DelMetaData(hash)
	assert.NoError(t, err, "DelMetaData should not return error")

}
