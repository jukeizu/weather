package geocoding

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/cache"
	"github.com/jukeizu/weather/api/protobuf-spec/geocodingpb"
)

type cacheServer struct {
	logger  log.Logger
	Service geocodingpb.GeocodeServer
	Cache   cache.Cache
}

func NewCacheServer(logger log.Logger, s geocodingpb.GeocodeServer, config cache.Config) geocodingpb.GeocodeServer {
	cache := cache.New(config)

	return &cacheServer{logger, s, cache}
}

func (s cacheServer) Geocode(ctx context.Context, req *geocodingpb.GeocodeRequest) (reply *geocodingpb.GeocodeReply, err error) {
	cacheResult := geocodingpb.GeocodeReply{}

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
