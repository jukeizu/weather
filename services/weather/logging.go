package weather

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	pb "github.com/jukeizu/weather/api/weather"
)

type loggingService struct {
	logger  log.Logger
	Service pb.WeatherServer
}

func NewLoggingService(logger log.Logger, s pb.WeatherServer) pb.WeatherServer {
	return &loggingService{logger, s}
}

func (s loggingService) Weather(ctx context.Context, req *pb.WeatherRequest) (reply *pb.WeatherReply, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Weather",
			"request", *req,
			"reply", *reply,
			"error", err,
			"took", time.Since(begin),
		)

	}(time.Now())

	reply, err = s.Service.Weather(ctx, req)

	return
}
