package repository

import (
	"DisHub/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOSSMetaRepository(t *testing.T) {
	// 初始化测试用的 OSSMetaRepository
	Repository := NewOSSMetaRepository(common.TEST_MYSQL_ADDR)
	hash := uuid.New().String()
	om := &OSSMetadata{
		Hash: hash,
		Size: 1024,
	}
	// 测试 PutMetaData 方法
	err := Repository.Insert(om)
	assert.NoError(t, err, "PutMetaData should not return error")
	err = Repository.Insert(om)
	assert.Error(t, err, "PutMetaData should return error")

	// 测试 GetMetaData 方法
	meta, err := Repository.GetMetaData(hash)
	assert.NoError(t, err, "GetMetaData should not return error")
	assert.NotNil(t, meta, "GetMetaData should return metadata")

	// 测试 DelMetaData 方法
	err = Repository.DelMetaData(hash)
	assert.NoError(t, err, "DelMetaData should not return error")

	// 关闭数据库连接
	Repository.Close()
}
