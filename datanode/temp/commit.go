package temp

import (
	"DisHub/common/utils"
	"DisHub/config"
	"DisHub/datanode/locate"
	"log"
	"os"
	"strconv"
	"strings"
)

func (t *tempInfo) hashAndId() (string, int) {
	s := strings.Split(t.Name, ".")
	if len(s) != 2 {
		return "", -1
	}
	id, _ := strconv.Atoi(s[1])
	return s[0], id
}

func commitTempObject(datFile string, tempinfo *tempInfo) {
	f, _ := os.Open(datFile)
	d := utils.CalculateHash(f)
	f.Close()
	targetPath := config.GetBasePath() + "/objects/" + tempinfo.Name + "." + d
	err := os.Rename(datFile, targetPath)
	if err != nil {
		log.Println(err)
	}
	tempinfo.hashAndId()
	locate.Add(tempinfo.hashAndId())
}
