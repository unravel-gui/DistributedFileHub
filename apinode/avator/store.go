package avator

import (
	"DisHub/common/utils"
	"fmt"
	"io"
	"net/http"
)

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	stream, e := putStream(hash, size)
	if e != nil {
		return http.StatusServiceUnavailable, e
	}
	reader := io.TeeReader(r, stream)
	d := utils.CalculateHash(reader)
	if d != hash {
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, calc=%s, request=%s", d, hash)
	}
	stream.Commit(true)
	return http.StatusOK, nil
}
