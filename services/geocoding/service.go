package geocoding

import (
	"context"
	"errors"

	pb "github.com/jukeizu/weather/api/geocoding"
	"googlemaps.github.io/maps"
)

type service struct {
	Client *maps.Client
}

func NewService(client *maps.Client) pb.GeocodeServer {
	return &service{client}
}

func (s service) Geocode(ctx context.Context, req *pb.GeocodeRequest) (*pb.GeocodeReply, error) {
	geocodingRequest := maps.GeocodingRequest{
		Address: req.Location,
	}

	results, err := s.Client.Geocode(context.Background(), &geocodingRequest)
	if err != nil {
		return nil, err
	}
	if len(results) < 1 {
		return nil, errors.New("No results for " + req.Location)

	}

	result := results[0]

	response := &pb.GeocodeReply{
		Latitude:         result.Geometry.Location.Lat,
		Longitude:        result.Geometry.Location.Lng,
		FormattedAddress: result.FormattedAddress,
	}

	return response, nil
}
