package weather

import (
	"context"
	"errors"

	"github.com/jukeizu/weather/api/protobuf-spec/geocodingpb"
	"github.com/jukeizu/weather/api/protobuf-spec/weatherpb"
	"github.com/shawntoffel/darksky"
)

type server struct {
	DarkSky       darksky.DarkSky
	GeocodeClient geocodingpb.GeocodeClient
}

func NewServer(darkskyClient darksky.DarkSky, geocodeClient geocodingpb.GeocodeClient) weatherpb.WeatherServer {
	return &server{darkskyClient, geocodeClient}
}

func (s server) Weather(ctx context.Context, req *weatherpb.WeatherRequest) (*weatherpb.WeatherReply, error) {
	geocodeRequest := &geocodingpb.GeocodeRequest{
		Location: req.Location,
	}

	location, err := s.GeocodeClient.Geocode(context.Background(), geocodeRequest)
	if err != nil {
		return nil, errors.New("geocode client error: " + err.Error())
	}

	forecastRequest := darksky.ForecastRequest{
		Time:      darksky.Timestamp(req.Time),
		Latitude:  darksky.Measurement(location.Latitude),
		Longitude: darksky.Measurement(location.Longitude),
	}

	if unitsAreValid(req.Units) {
		forecastRequest.Options.Units = req.Units
	}

	darkskyResponse, err := s.DarkSky.Forecast(forecastRequest)
	if err != nil {
		return nil, err
	}

	mapper := NewMapper(darkskyResponse)

	return mapper.AsWeatherResponse(location)
}

func unitsAreValid(units string) bool {
	return units == "us" || units == "si"
}
