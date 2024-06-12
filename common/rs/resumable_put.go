package rs

import (
	"DisHub/common/objectstream"
	"DisHub/common/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type resumableToken struct {
	Size    int64
	Hash    string
	Servers []string
	Uuids   []string
}

type RSResumablePutStream struct {
	*RSPutStream
	*resumableToken
}

func NewRSResumablePutStream(dataServers []string, hash string, size int64) (*RSResumablePutStream, error) {
	putStream, e := NewRSPutStream(dataServers, hash, size)
	if e != nil {
		return nil, e
	}
	uuids := make([]string, ALL_SHARDS)
	for i := range uuids {
		uuids[i] = putStream.writers[i].(*objectstream.TempPutStream).Uuid
	}
	token := &resumableToken{size, hash, dataServers, uuids}
	return &RSResumablePutStream{putStream, token}, nil
}

func NewRSResumablePutStreamFromToken(token string) (*RSResumablePutStream, error) {
	// token 即序列化后的resumableToken
	b, e := base64.URLEncoding.DecodeString(token)
	if e != nil {
		return nil, e
	}

	var t resumableToken
	e = json.Unmarshal(b, &t)
	if e != nil {
		return nil, e
	}

	writers := make([]io.Writer, ALL_SHARDS)
	for i := range writers {
		// 实现了write方法即为io.writer
		writers[i] = &objectstream.TempPutStream{t.Servers[i], t.Uuids[i]}
	}
	enc := NewEncoder(writers)
	return &RSResumablePutStream{&RSPutStream{enc}, &t}, nil
}

func (s *RSResumablePutStream) ToToken() string {
	b, _ := json.Marshal(s)
	return base64.URLEncoding.EncodeToString(b)
}

// CurrentSize 获得文件大小
func (s *RSResumablePutStream) CurrentSize() int64 {
	r, e := http.Head(fmt.Sprintf("http://%s/temp/%s", s.Servers[0], s.Uuids[0]))
	if e != nil {
		log.Println(e)
		return -1
	}
	if r.StatusCode != http.StatusOK {
		log.Println(r.StatusCode)
		return -1
	}
	// 获得切片文件的大小，乘保存数据的节点即为文件大小，如果大于只是代表文件不是被整数划分而已
	size := utils.GetSizeFromHttpHeader(r.Header) * DATA_SHARDS
	if size > s.Size {
		size = s.Size
	}
	return size
}
