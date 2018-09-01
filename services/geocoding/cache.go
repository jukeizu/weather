package geocoding

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	pb "github.com/jukeizu/weather/api/geocoding"
	"github.com/shawntoffel/services-core/cache"
)

type cacheService struct {
	logger  log.Logger
	Service pb.GeocodeServer
	Cache   cache.Cache
}

func NewCacheService(logger log.Logger, s pb.GeocodeServer, config cache.Config) pb.GeocodeServer {
	cache := cache.NewCache(config)

	return &cacheService{logger, s, cache}
}

func (s cacheService) Geocode(ctx context.Context, req *pb.GeocodeRequest) (reply *pb.GeocodeReply, err error) {
	cacheResult := pb.GeocodeReply{}

	cacheErr := s.Cache.Get(req, &cacheResult)
	if cacheErr == nil {
		return &cacheResult, nil
	}

	reply, err = s.Service.Geocode(ctx, req)
	if err != nil {
		return reply, err
	}

	err = s.Cache.Set(req, reply, time.Hour*480)

	return
}
