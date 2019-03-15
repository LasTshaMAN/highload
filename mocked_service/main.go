package main

import (
	"highload/http_adapters"
	"highload/mocked_service/api"
	"highload/mocked_service/api/adapters/iris"
	"highload/mocked_service/domain"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	v := domain.NewValuer()
	s := domain.NewSleeper()
	a := api.New(v, s)
	i := iris.New(a)
	// TODO
	// make port dynamic
	if _, err := httpAdapters.RunIris(i, 8002); err != nil {
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
