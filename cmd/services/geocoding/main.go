package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/jukeizu/weather/api/geocoding"
	"github.com/jukeizu/weather/services/geocoding"
	"github.com/shawntoffel/services-core/cache"
	"github.com/shawntoffel/services-core/command"
	"github.com/shawntoffel/services-core/config"
	"github.com/shawntoffel/services-core/logging"
	"google.golang.org/grpc"
	"googlemaps.github.io/maps"
)

var serviceArgs command.CommandArgs

func init() {
	serviceArgs = command.ParseArgs()
}

type Config struct {
	Port        int
	ApiKey      string
	CacheConfig cache.Config
}

func main() {
	logger := logging.GetLogger("services.geocoding", os.Stdout)

	c := Config{}
	err := config.ReadConfig(serviceArgs.ConfigFile, &c)
	if err != nil {
		panic(err)
	}

	client, err := maps.NewClient(maps.WithAPIKey(c.ApiKey))
	if err != nil {
		panic(err)
	}

	service := geocoding.NewService(client)
	service = geocoding.NewCacheService(logger, service, c.CacheConfig)
	service = geocoding.NewLoggingService(logger, service)

	errChannel := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errChannel <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		port := fmt.Sprintf(":%d", c.Port)

		listener, err := net.Listen("tcp", port)
		if err != nil {
			logger.Log("error", err.Error())
		}

		s := grpc.NewServer()
		pb.RegisterGeocodeServer(s, service)

		logger.Log("transport", "grpc", "address", port, "msg", "listening")

		errChannel <- s.Serve(listener)
	}()

	logger.Log("stopped", <-errChannel)
}
