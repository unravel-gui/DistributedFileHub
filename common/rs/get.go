package rs

import (
	"DisHub/common/objectstream"
	"fmt"
	"io"
)

type RSGetStream struct {
	*decoder
}

func NewRSGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*RSGetStream, error) {
	// 节点数量对不上
	if len(locateInfo)+len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}
	// 响应节点在前，未响应节点在后
	readers := make([]io.Reader, ALL_SHARDS)
	for i := 0; i < ALL_SHARDS; i++ {
		server := locateInfo[i]
		// 将未响应的ds节点补充到后面
		if server == "" {
			locateInfo[i] = dataServers[0]
			dataServers = dataServers[1:]
			continue
		}
		// 获得块文件内容
		reader, e := objectstream.NewGetStream(server, fmt.Sprintf("%s.%d", hash, i))
		if e == nil {
			readers[i] = reader
		}
	}

	writers := make([]io.Writer, ALL_SHARDS)
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	var e error
	// 创建需要恢复数据的任务
	for i := range readers {
		if readers[i] == nil {
			// 为空说明数据丢失
			// 恢复数据,创建上传任务,调用write方法会将数据上传到ds节点
			writers[i], e = objectstream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", hash, i), perShard)
			// 创建任务失败则报错返回
			if e != nil {
				return nil, e
			}
		}
	}
	// 封装读取任务和恢复任务
	dec := NewDecoder(readers, writers, size)
	return &RSGetStream{dec}, nil
}

func (s *RSGetStream) Close() {
	for i := range s.writers {
		if s.writers[i] != nil {
			s.writers[i].(*objectstream.TempPutStream).Commit(true)
		}
	}
}

// Seek 移动指针
func (s *RSGetStream) Seek(offset int64, whence int) (int64, error) {
	if whence != io.SeekCurrent {
		return 0, fmt.Errorf("only support SeekCurrent")
	}
	if offset < 0 {
		return 0, fmt.Errorf("only support forward seek")
	}
	for offset != 0 {
		length := int64(BLOCK_SIZE)
		if offset < length {
			length = offset
		}
		buf := make([]byte, length)
		io.ReadFull(s, buf)
		offset -= length
	}
	return offset, nil
}
