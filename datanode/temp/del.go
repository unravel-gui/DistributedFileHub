package temp

import (
	"DisHub/config"
	"github.com/gin-gonic/gin"
	"os"
)

func del(c *gin.Context) {
	uuid := c.Param("uuid")
	infoFile := config.GetBasePath() + "/temp/" + uuid
	datFile := infoFile + ".dat"
	os.Remove(infoFile)
	os.Remove(datFile)
}
