package rs

import (
	"github.com/klauspost/reedsolomon"
	"io"
)

type encoder struct {
	writers []io.Writer
	enc     reedsolomon.Encoder
	cache   []byte
}

func NewEncoder(writers []io.Writer) *encoder {
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &encoder{writers, enc, nil}
}

func (e *encoder) Write(p []byte) (n int, err error) {
	length := len(p)
	current := 0
	for length != 0 {
		// 计算当前能放多少
		next := BLOCK_SIZE - len(e.cache)
		// 当前块放不满
		if next > length {
			next = length
		}
		// 数据放入缓冲区
		e.cache = append(e.cache, p[current:current+next]...)
		// 当前块放满了
		if len(e.cache) == BLOCK_SIZE {
			e.Flush()
		}
		current += next
		length -= next
	}
	return len(p), nil
}

func (e *encoder) Flush() {
	if len(e.cache) == 0 {
		return
	}
	// 分割当前块文件,纠删码的分片
	shards, _ := e.enc.Split(e.cache)
	// 对分片后的数据进行编码
	e.enc.Encode(shards)
	// 分片写入切片节点,对应
	for i := range shards {
		// 实际调用TempPutStream，将数据上传到datanode中
		e.writers[i].Write(shards[i])
	}
	e.cache = []byte{}
}
