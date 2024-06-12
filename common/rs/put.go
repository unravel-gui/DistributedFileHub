package rs

import (
	"DisHub/common/objectstream"
	"fmt"
	"io"
)

type RSPutStream struct {
	*encoder
}

func NewRSPutStream(dataServers []string, hash string, size int64) (*RSPutStream, error) {
	// 校验ds个数
	if len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}
	// 保证存储空间大于存储内容
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	writers := make([]io.Writer, ALL_SHARDS)
	var e error
	// 为每个分片节点创建写入流
	for i := range writers {
		// 文件名称为hash.id
		writers[i], e = objectstream.NewTempPutStream(dataServers[i],
			fmt.Sprintf("%s.%d", hash, i), perShard)
		if e != nil {
			return nil, e
		}
	}
	enc := NewEncoder(writers)

	return &RSPutStream{enc}, nil
}

func (s *RSPutStream) Commit(success bool) {
	// 将不满足整块的文件写入writer中
	s.Flush()
	for i := range s.writers {
		s.writers[i].(*objectstream.TempPutStream).Commit(success)
	}
}
