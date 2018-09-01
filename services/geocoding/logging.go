package geocoding

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	pb "github.com/jukeizu/weather/api/geocoding"
)

type loggingService struct {
	logger  log.Logger
	Service pb.GeocodeServer
}

func NewLoggingService(logger log.Logger, s pb.GeocodeServer) pb.GeocodeServer {
	return &loggingService{logger, s}
}

func (s loggingService) Geocode(ctx context.Context, req *pb.GeocodeRequest) (reply *pb.GeocodeReply, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Geocode",
			"request", *req,
			"reply", *reply,
			"error", err,
			"took", time.Since(begin),
		)

	}(time.Now())

	reply, err = s.Service.Geocode(ctx, req)

	return
}
