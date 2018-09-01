package weather

import (
	"context"

	gpb "github.com/jukeizu/weather/api/geocoding"
	wpb "github.com/jukeizu/weather/api/weather"
	"github.com/shawntoffel/darksky"
)

type service struct {
	DarkSky       darksky.DarkSky
	GeocodeClient gpb.GeocodeClient
}

func NewService(darkskyClient darksky.DarkSky, geocodeClient gpb.GeocodeClient) wpb.WeatherServer {
	return &service{darkskyClient, geocodeClient}
}

func (s service) Weather(ctx context.Context, req *wpb.WeatherRequest) (*wpb.WeatherReply, error) {
	geocodeRequest := &gpb.GeocodeRequest{
		Location: req.Location,
	}

	location, err := s.GeocodeClient.Geocode(context.Background(), geocodeRequest)
	if err != nil {
		return nil, err
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
