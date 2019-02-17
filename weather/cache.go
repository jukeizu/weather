package weather

import (
	"context"
	"time"

	"github.com/jukeizu/cache"
	"github.com/jukeizu/weather/api/protobuf-spec/weatherpb"
	"github.com/rs/zerolog"
)

type cacheServer struct {
	logger zerolog.Logger
	Server weatherpb.WeatherServer
	Cache  cache.Cache
}

func NewCacheServer(logger zerolog.Logger, s weatherpb.WeatherServer, config cache.Config) weatherpb.WeatherServer {
	cache := cache.New(config)

	return &cacheServer{logger, s, cache}
}

func (s cacheServer) Weather(ctx context.Context, req *weatherpb.WeatherRequest) (*weatherpb.WeatherReply, error) {
	cacheResult := weatherpb.WeatherReply{}

	cacheErr := s.Cache.Get(req, &cacheResult)
	if cacheErr == nil {
		s.logger.Debug().Msg("found cached reply")
		return &cacheResult, nil
	}

	s.logger.Debug().Err(cacheErr).Msg("could not fetch from cache")

	reply, err := s.Server.Weather(ctx, req)
	if err != nil {
		return reply, err
	}

	err = s.Cache.Set(req, reply, time.Minute*20)
	if err != nil {
		s.logger.Debug().Err(err).Msg("could not set cache")
	}

	return reply, nil
}
