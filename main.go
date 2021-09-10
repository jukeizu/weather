package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpczerolog "github.com/cheapRoc/grpc-zerolog"
	_ "github.com/jnewmano/grpc-json-proxy/codec"
	"github.com/jukeizu/cache"
	"github.com/jukeizu/weather/api/protobuf-spec/geocodingpb"
	"github.com/jukeizu/weather/api/protobuf-spec/weatherpb"
	"github.com/jukeizu/weather/geocoding"
	"github.com/jukeizu/weather/treediagram"
	"github.com/jukeizu/weather/weather"
	"github.com/oklog/run"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	azweather "github.com/shawntoffel/azure-maps-go/weather"
	"github.com/shawntoffel/darksky"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
	"googlemaps.github.io/maps"
)

var Version = ""

var (
	flagVersion = false
	flagDebug   = false
	flagServer  = false
	flagHandler = false

	grpcPort       = "50052"
	httpPort       = "10002"
	serviceAddress = "localhost:" + grpcPort
	cacheAddress   = cache.DefaultRedisAddress
)

func init() {
	flag.StringVar(&grpcPort, "grpc.port", grpcPort, "grpc port for server")
	flag.StringVar(&httpPort, "http.port", httpPort, "http port for handler")
	flag.StringVar(&cacheAddress, "cache.addr", cacheAddress, "cache address")
	flag.StringVar(&serviceAddress, "service.addr", serviceAddress, "sercice address if not local")
	flag.BoolVar(&flagServer, "server", false, "Run as server")
	flag.BoolVar(&flagHandler, "handler", false, "Run as handler")
	flag.BoolVar(&flagVersion, "v", false, "version")
	flag.BoolVar(&flagDebug, "D", false, "enable debug logging")

	flag.Parse()
}

func main() {
	if flagVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().
		Str("instance", xid.New().String()).
		Str("component", "weather").
		Str("version", Version).
		Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if flagDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	grpcLoggerV2 := grpczerolog.New(logger.With().Str("transport", "grpc").Logger())
	grpclog.SetLoggerV2(grpcLoggerV2)

	if !flagServer && !flagHandler {
		flagServer = true
		flagHandler = true
	}

	clientConn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                30 * time.Second,
				Timeout:             10 * time.Second,
				PermitWithoutStream: true,
			},
		),
	)
	if err != nil {
		logger.Error().Err(err).Str("serviceAddress", serviceAddress).Msg("could not dial service address")
		os.Exit(1)
	}

	g := run.Group{}

	if flagServer {
		grpcServer := newGrpcServer(logger)
		server := NewServer(logger, grpcServer)

		darkskyTokenFile := os.Getenv("DARKSKY_TOKEN_FILE")
		darkskyToken := readSecretsFile(logger, darkskyTokenFile)
		darkskyClient := darksky.New(darkskyToken)

		azmapsTokenFile := os.Getenv("AZURE_MAPS_TOKEN_FILE")
		azmapsToken := readSecretsFile(logger, azmapsTokenFile)
		azweatherClient := azweather.New(azmapsToken)

		mapsTokenFile := os.Getenv("GOOGLE_MAPS_TOKEN_FILE")
		mapsToken := readSecretsFile(logger, mapsTokenFile)
		mapsClient, err := maps.NewClient(maps.WithAPIKey(mapsToken))
		if err != nil {
			logger.Error().Err(err).Msg("could not start maps client")
			os.Exit(1)
		}

		cacheConfig := cache.Config{
			Address: cacheAddress,
			Version: Version,
		}

		geocodingServer := geocoding.NewServer(mapsClient)
		geocodingServer = geocoding.NewCacheServer(logger.With().Str("component", "geocoding").Logger(), geocodingServer, cacheConfig)
		geocodingpb.RegisterGeocodeServer(grpcServer, geocodingServer)

		geocodeClient := geocodingpb.NewGeocodeClient(clientConn)
		weatherServer := weather.NewServer(darkskyClient, azweatherClient, geocodeClient)
		weatherServer = weather.NewCacheServer(logger, weatherServer, cacheConfig)
		weatherpb.RegisterWeatherServer(grpcServer, weatherServer)

		grpcAddr := ":" + grpcPort

		g.Add(func() error {
			return server.Start(grpcAddr)
		}, func(error) {
			server.Stop()
		})
	}

	if flagHandler {
		client := weatherpb.NewWeatherClient(clientConn)
		httpAddr := ":" + httpPort

		handler := treediagram.NewHandler(logger, client, httpAddr)

		g.Add(func() error {
			return handler.Start()
		}, func(error) {
			err := handler.Stop()
			if err != nil {
				logger.Error().Err(err).Caller().Msg("couldn't stop handler")
			}
		})
	}

	cancel := make(chan struct{})
	g.Add(func() error {
		return interrupt(cancel)
	}, func(error) {
		close(cancel)
	})

	logger.Info().Err(g.Run()).Msg("stopped")
}

func newGrpcServer(logger zerolog.Logger) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    5 * time.Minute,
				Timeout: 10 * time.Second,
			},
		),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             5 * time.Second,
				PermitWithoutStream: true,
			},
		),
		LoggingInterceptor(logger),
	)

	return grpcServer
}

func readSecretsFile(logger zerolog.Logger, filename string) string {
	tokenBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error().Err(err).
			Str("filename", filename).
			Msg("could not read secrets file")
		os.Exit(1)
	}

	return string(tokenBytes)
}

func interrupt(cancel <-chan struct{}) error {
	c := make(chan os.Signal, 0)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-cancel:
		return errors.New("stopping")
	case sig := <-c:
		return fmt.Errorf("%s", sig)
	}
}
