package main

import (
	"bitbucket.org/victorcoder/dcron/dcron"
)

func main() {
	s := dcron.NewScheduler()
	s.Load()
	dcron.InitSerfAgent()
	dcron.ServerInit()
}
