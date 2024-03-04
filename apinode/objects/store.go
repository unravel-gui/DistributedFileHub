package objects

import (
	"DisHub/apinode/locate"
	"DisHub/common/utils"
	"fmt"
	"io"
	"net/http"
)

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	// 存在,不重复上传
	if locate.Exist(hash) {
		return http.StatusOK, nil
	}
	stream, e := putStream(hash, size)
	if e != nil {
		return http.StatusServiceUnavailable, e
	}
	reader := io.TeeReader(r, stream)
	d := utils.CalculateHash(reader)
	if d != hash {
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, calc=%d, request=%d", d, hash)
	}
	stream.Commit(true)
	return http.StatusOK, nil
}
