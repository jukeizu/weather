package geocoding

import (
	"context"
	"errors"

	"github.com/jukeizu/weather/api/protobuf-spec/geocodingpb"
	"googlemaps.github.io/maps"
)

type server struct {
	Client *maps.Client
}

func NewServer(client *maps.Client) geocodingpb.GeocodeServer {
	return &server{client}
}

func (s server) Geocode(ctx context.Context, req *geocodingpb.GeocodeRequest) (*geocodingpb.GeocodeReply, error) {
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

	response := &geocodingpb.GeocodeReply{
		Latitude:         result.Geometry.Location.Lat,
		Longitude:        result.Geometry.Location.Lng,
		FormattedAddress: result.FormattedAddress,
	}

	return response, nil
}
