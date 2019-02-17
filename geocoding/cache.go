package geocoding

import (
	"context"
	"time"

	"github.com/jukeizu/cache"
	"github.com/jukeizu/weather/api/protobuf-spec/geocodingpb"
	"github.com/rs/zerolog"
)

type cacheServer struct {
	logger zerolog.Logger
	Server geocodingpb.GeocodeServer
	Cache  cache.Cache
}

func NewCacheServer(logger zerolog.Logger, s geocodingpb.GeocodeServer, config cache.Config) geocodingpb.GeocodeServer {
	cache := cache.New(config)

	return &cacheServer{logger, s, cache}
}

func (s cacheServer) Geocode(ctx context.Context, req *geocodingpb.GeocodeRequest) (*geocodingpb.GeocodeReply, error) {
	cacheResult := geocodingpb.GeocodeReply{}

	cacheErr := s.Cache.Get(req, &cacheResult)
	if cacheErr == nil {
		s.logger.Debug().Msg("found cached reply")
		return &cacheResult, nil
	}

	s.logger.Debug().Err(cacheErr).Msg("could not fetch from cache")

	reply, err := s.Server.Geocode(ctx, req)
	if err != nil {
		return reply, err
	}

	err = s.Cache.Set(req, reply, time.Hour*480)
	if err != nil {
		s.logger.Debug().Err(err).Msg("could not set cache")
	}

	return reply, nil
}
