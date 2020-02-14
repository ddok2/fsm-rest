package server

import (
	"io/ioutil"
	"net/http"

	"blockchain.automation/model"
	"blockchain.automation/utils"
	"github.com/gin-gonic/gin"
)

var (
	logger *utils.Logger
)

func pause(c *gin.Context) {

	logger.Info("pause")

	model.SetExit(true)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "OK"})

}

func resume(c *gin.Context) {

	logger.Info("resume")

	model.Resume()
	model.SetExit(false)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "OK"})

}

func report(c *gin.Context) {

	logger.Info("report")

	members := model.NewMembers()
	members.Report()
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "OK"})

}

func SetupRouter() *gin.Engine {

	logger = utils.NewLogger()

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	router := gin.Default()

	router.GET("/action/pause", pause)
	router.GET("/action/resume", resume)
	router.GET("/action/report", report)

	return router
}
