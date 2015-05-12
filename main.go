package main

import (
	"bitbucket.org/victorcoder/dcron/dcron"
	"time"
)

func main() {
	s := dcron.NewScheduler()
	s.Load()
	dcron.ServerInit()
	time.Sleep(2 * time.Minute)
}
