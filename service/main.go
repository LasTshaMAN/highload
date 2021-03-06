package main

import (
	"crypto/tls"
	"highload/http_adapters"
	"highload/service/api"
	"highload/service/api/iris"
	"highload/service/api/iris/middleware"
	"highload/service/domain"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/http2"

	"github.com/sirupsen/logrus"
)

const serviceName = "service"

func main() {
	client := &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		// TODO
		// Play with client and server timeouts
		Timeout: 300 * time.Millisecond,
	}
	avg := domain.NewConcurrentAvg("https://127.0.0.1:8002", client, 2*32*1024)
	a := api.New(avg)
	i := iris.New(a, middleware.NewPrometheus(serviceName))
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
