package stub

import (
	"DisHub/common/utils"
	"encoding/json"
	"net/http"
)

func GetDataServers(server string) []string {
	return GetServers(server, "datanode")
}
func GetUserServers(server string) []string {
	return GetServers(server, "usernode")
}

func GetServers(server, path string) []string {
	request, e := http.NewRequest("GET", "http://"+server+"/"+path, nil)
	if e != nil {
		return nil
	}
	utils.SetMagicTokenFromHeader(request.Header)
	client := http.Client{}
	resp, e := client.Do(request)
	if e != nil {
		return nil
	}
	var ds []string
	defer resp.Body.Close()
	// 读取和解码响应
	_ = json.NewDecoder(resp.Body).Decode(&ds)
	return ds
}

func GetGetLoadBalanceUserServer(server string) string {
	return GetServer(server, "usernode")
}

func GetServer(server, path string) string {
	request, e := http.NewRequest("GET", "http://"+server+"/lb/"+path, nil)
	if e != nil {
		return ""
	}
	utils.SetMagicTokenFromHeader(request.Header)
	client := http.Client{}
	resp, e := client.Do(request)
	if e != nil {
		return ""
	}
	var s string
	defer resp.Body.Close()
	// 读取和解码响应
	_ = json.NewDecoder(resp.Body).Decode(&s)
	return s
}
