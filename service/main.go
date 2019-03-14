package main

import (
	"highload/service/api"
	"highload/service/api/adapters/iris"
	"highload/service/domain"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	v := domain.NewValuer()
	s := domain.NewSleeper()
	a := api.New(v, s)
	// TODO
	// make port dynamic
	if _, err := iris.Run(a, 8000); err != nil {
		logrus.Error(err)
		return
	}

	waitForExitSignal()
}

func waitForExitSignal() {
	// Graceful shutdown can be implemented here ...

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	<-exit
}
