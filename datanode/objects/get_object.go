package objects

import (
	"DisHub/common/utils"
	"DisHub/config"
	"DisHub/datanode/locate"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getFile(name string) string {
	files, _ := filepath.Glob(config.GetBasePath() + "/objects/" + name + ".*")
	if len(files) != 1 {
		return ""
	}
	file := files[0]
	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return ""
	}
	d := utils.CalculateHash(f)
	hash := strings.Split(file, ".")[2]
	if d != hash {
		locate.Del(hash)
		os.Remove(file)
		return ""
	}
	return file
}
