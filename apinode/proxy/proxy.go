package proxy

import (
	"DisHub/apinode/heartbeat"
	"DisHub/common"
	"DisHub/common/response"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func userProxy(c *gin.Context) {
	ss, err := heartbeat.GetBalanceUserServer()
	if err != nil {
		response.InternalServerByError(c, err)
		return
	}
	if ss == "" {
		response.InternalServer(c, "no enough user node")
		return
	}
	ok := forward(c, ss)
	if ok {
		return
	}
	response.InternalServer(c, "forward user node failed")
}

func masterProxy(c *gin.Context) {
	ss := heartbeat.GetMasterServers()
	if len(ss) == 0 {
		response.InternalServer(c, "no enough master")
		return
	}
	common.G_Random.Shuffle(len(ss), func(i, j int) {
		ss[i], ss[j] = ss[j], ss[i]
	})
	for _, server := range ss {
		ok := forward(c, server)
		if ok {
			return
		}
	}
	response.InternalServer(c, "forward user node failed")
}

func forward(c *gin.Context, server string) bool {
	// 构建要转发的请求
	targetURL := "http://" + server + c.Request.URL.Path
	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		return false
	}
	// 将当前请求的请求头复制到转发请求中
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// 发送请求到目标服务器
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	// 将响应头写入到当前响应中
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Set(key, value)
		}
	}
	// 返回响应给客户端
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
	return true
}
