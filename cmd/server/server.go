package server

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/stkr89/go-auth-service-template/common"
	"github.com/stkr89/go-auth-service-template/endpoints"
	"github.com/stkr89/go-auth-service-template/service"
	transport "github.com/stkr89/go-auth-service-template/transports"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log/level"
	"github.com/joho/godotenv"
)

func InitServer() {
	logger := common.NewLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Log("message", ".env file not found", "err", err)
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	e := endpoints.MakeEndpoints(service.NewAuthServiceImpl())
	StartServer(logger, e, true, true)

	level.Error(logger).Log("exit", <-errs)
}

func StartServer(logger log.Logger, e endpoints.Endpoints, startGRPC, startHTTP bool) {
	if startHTTP {
		startHTTPServer(logger, e)
	}
}

func startHTTPServer(logger log.Logger, e endpoints.Endpoints) {
	httpPort := os.Getenv("HTTP_PORT")
	listener, err := getListener(httpPort)
	if err != nil {
		logger.Log("transport", "HTTP", "during", "Listen", "err", err)
		os.Exit(1)
	}

	httpHandler := transport.NewHTTPHandler(e)

	go func() {
		level.Info(logger).Log("msg", fmt.Sprintf("Starting HTTP server 🚀 at %s", httpPort))
		http.Serve(listener, httpHandler)
	}()
}

func getListener(port string) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf(":%s", port))
}
