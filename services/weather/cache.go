package weather

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	pb "github.com/jukeizu/weather/api/weather"
	"github.com/shawntoffel/services-core/cache"
)

type cacheService struct {
	logger  log.Logger
	Service pb.WeatherServer
	Cache   cache.Cache
}

func NewCacheService(logger log.Logger, s pb.WeatherServer, config cache.Config) pb.WeatherServer {
	cache := cache.NewCache(config)

	return &cacheService{logger, s, cache}
}

func (s cacheService) Weather(ctx context.Context, req *pb.WeatherRequest) (reply *pb.WeatherReply, err error) {
	cacheResult := pb.WeatherReply{}

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
