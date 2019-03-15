package main

import (
	"highload/http_adapters"
	"highload/service/api"
	"highload/service/api/iris"
	"highload/service/domain"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	avg := domain.NewAvg()
	a := api.New(avg)
	i := iris.New(a)
	// TODO
	// make port dynamic
	if _, err := httpAdapters.RunIris(i, 8001); err != nil {
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
