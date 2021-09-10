package weather

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jukeizu/weather/api/protobuf-spec/geocodingpb"
	"github.com/jukeizu/weather/api/protobuf-spec/weatherpb"
	"github.com/jukeizu/weather/weather/azmapssource"
	"github.com/jukeizu/weather/weather/darkskysource"
	azweather "github.com/shawntoffel/azure-maps-go/weather"
	"github.com/shawntoffel/darksky"
)

type server struct {
	DarkSky       darksky.DarkSky
	AzMaps        azweather.Weather
	GeocodeClient geocodingpb.GeocodeClient
}

func NewServer(
	darkskyClient darksky.DarkSky,
	azMapsClient azweather.Weather,
	geocodeClient geocodingpb.GeocodeClient,
) weatherpb.WeatherServer {
	return &server{
		darkskyClient,
		azMapsClient,
		geocodeClient,
	}
}

func (s server) Weather(ctx context.Context, req *weatherpb.WeatherRequest) (*weatherpb.WeatherReply, error) {
	geocodeRequest := &geocodingpb.GeocodeRequest{
		Location: req.Location,
	}

	location, err := s.GeocodeClient.Geocode(context.Background(), geocodeRequest)
	if err != nil {
		return nil, errors.New("geocode client error: " + err.Error())
	}

	if req.Source == "" || strings.EqualFold(req.Source, "darksky") {
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

		mapper := darkskysource.NewMapper(darkskyResponse)

		return mapper.AsWeatherResponse(location)
	}

	if strings.EqualFold(req.Source, "azure") {
		opts := &azweather.CurrentConditionsRequestOptions{
			Unit: req.Units,
		}
		azResponse, err := s.AzMaps.CurrentConditions(fmt.Sprintf("%f,%f", location.Latitude, location.Longitude), opts)
		if err != nil {
			return nil, err
		}

		mapper := azmapssource.NewMapper(azResponse.Results[0])

		return mapper.AsWeatherResponse(location)
	}

	return nil, errors.New("unknown weather source")
}

func unitsAreValid(units string) bool {
	return units == "us" || units == "si"
}
