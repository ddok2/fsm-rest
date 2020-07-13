package main

import (
	"runtime"

	cfg "blockchain.automation/common/config"
	"blockchain.automation/core/bot/action"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	config := cfg.NewConfig()
	if err := config.Load("./common/config/config.yml"); err != nil {
		return
	}

	act := action.NewAction(config)

	adminBot := act.GetAdminBot()
	adminBot.Balance = 219996170000

	act.Handler()
}
