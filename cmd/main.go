package main

import (
	"runtime"

	"blockchain.automation/cmd/common"
	"blockchain.automation/cmd/server"
	"blockchain.automation/fsm"
	"blockchain.automation/model"
	"blockchain.automation/utils"
)

var (
	logger  *utils.Logger
	config  *common.Config
	members *model.Members
)

var events fsm.Events

func initialize() error {

	logger = utils.NewLogger()
	logger.InitLogger()

	config = common.NewConfig()

	config.LoadConfigYaml()

	members = model.NewMembers()

	members.Initialize()
	members.CreateMembers()
	if config.OperationMode == "blockchain" {
		members.RegisterMembers()
	} else {
		members.Signup()
	}

	return nil
}

func finalize() {
	logger.Finalize()
}

func run() {
	router := server.SetupRouter()

	go func() {
		if err := router.Run(":" + "8090"); err != nil {
			logger.Error(err.Error())
			return
		}
	}()

}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := initialize(); err != nil {
		logger.Error(err.Error())
		return
	}

	defer finalize()

	go run()

	go members.Schedule()

	end := make(chan string)

	<-end

}
