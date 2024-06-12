package usernode

import (
	"DisHub/master/heartbeat"
	"github.com/gin-gonic/gin"
	"net/http"
)

func get(c *gin.Context) {
	ds := heartbeat.GetUserServers()
	c.JSON(http.StatusOK, ds)
}
