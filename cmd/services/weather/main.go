package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	gpb "github.com/jukeizu/weather/api/geocoding"
	wpb "github.com/jukeizu/weather/api/weather"
	"github.com/jukeizu/weather/services/weather"
	"github.com/shawntoffel/darksky"
	"github.com/shawntoffel/services-core/cache"
	"github.com/shawntoffel/services-core/command"
	"github.com/shawntoffel/services-core/config"
	"github.com/shawntoffel/services-core/logging"
	"google.golang.org/grpc"
)

var serviceArgs command.CommandArgs

func init() {
	serviceArgs = command.ParseArgs()
}

type Config struct {
	Port            int
	DarkSkyApiKey   string
	CacheConfig     cache.Config
	GeocodeEndpoint string
}

func main() {
	logger := logging.GetLogger("services.weather", os.Stdout)

	c := Config{}
	err := config.ReadConfig(serviceArgs.ConfigFile, &c)
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial(c.GeocodeEndpoint, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	geocodeClient := gpb.NewGeocodeClient(conn)

	darkskyClient := darksky.New(c.DarkSkyApiKey)

	service := weather.NewService(darkskyClient, geocodeClient)
	service = weather.NewCacheService(logger, service, c.CacheConfig)
	service = weather.NewLoggingService(logger, service)

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
		wpb.RegisterWeatherServer(s, service)

		logger.Log("transport", "grpc", "address", port, "msg", "listening")

		errChannel <- s.Serve(listener)
	}()

	logger.Log("stopped", <-errChannel)
}
