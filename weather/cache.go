package weather

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/cache"
	"github.com/jukeizu/weather/api/protobuf-spec/weatherpb"
)

type cacheService struct {
	logger  log.Logger
	Service weatherpb.WeatherServer
	Cache   cache.Cache
}

func NewCacheService(logger log.Logger, s weatherpb.WeatherServer, config cache.Config) weatherpb.WeatherServer {
	cache := cache.New(config)

	return &cacheService{logger, s, cache}
}

func (s cacheService) Weather(ctx context.Context, req *weatherpb.WeatherRequest) (reply *weatherpb.WeatherReply, err error) {
	cacheResult := weatherpb.WeatherReply{}

	cacheErr := s.Cache.Get(req, &cacheResult)
	if cacheErr == nil {
		return &cacheResult, nil
	}

	reply, err = s.Service.Weather(ctx, req)
	if err != nil {
		return reply, err
	}

	err = s.Cache.Set(req, reply, time.Minute*20)

	return
}
