package temp

import (
	"DisHub/config"
	"DisHub/datanode/locate"
	"log"
	"os"
)

func commitTempObject(datFile string, tempinfo *tempInfo) {
	targetPath := config.GetBasePath() + "/objects/" + tempinfo.Name
	err := os.Rename(datFile, targetPath)
	if err != nil {
		log.Println(err)
	}
	locate.Add(tempinfo.Name)
}
