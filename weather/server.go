package weather

import (
	"context"
	"errors"
	"fmt"

	"github.com/jukeizu/weather/api/protobuf-spec/geocodingpb"
	"github.com/jukeizu/weather/api/protobuf-spec/weatherpb"
	"github.com/jukeizu/weather/weather/azsource"
	"github.com/jukeizu/weather/weather/darkskysource"
	"github.com/shawntoffel/azure-maps-go/azweather"
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
	location, err := s.lookupLocation(req.Location)
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

	mapper := darkskysource.NewMapper(darkskyResponse)

	return mapper.AsWeatherResponse(location)
}

func (s server) Plan(ctx context.Context, req *weatherpb.PlanRequest) (*weatherpb.PlanReply, error) {
	location, err := s.lookupLocation(req.Location)
	if err != nil {
		return nil, err
	}

	duration := int(req.Duration)
	if duration == 0 {
		duration = 120
	}

	units := req.Units
	if units == "" {
		units = "imperial"
	}

	opts := azweather.HourlyForecastRequestOptions{
		Duration: &duration,
		Unit:     units,
	}

	query := fmt.Sprintf("%g,%g", location.Latitude, location.Longitude)

	azResponse, err := s.AzMaps.HourlyForecast(query, &opts)
	if err != nil {
		return nil, err
	}

	hours := []azweather.HourlyForecast{}

	for _, hour := range azResponse.Forecasts {
		if req.Daylight && !hour.IsDaylight {
			continue
		}

		if !isValidDataRange(req.Wind, hour.Wind.Speed.Value) {
			continue
		}

		if !isValidDataRange(req.WindGust, hour.WindGust.Speed.Value) {
			continue
		}

		if !isValidDataRange(req.Temperature, hour.Temperature.Value) {
			continue
		}

		if req.Precipitation && !hour.HasPrecipitation {
			continue
		}

		hours = append(hours, hour)
	}

	mapper := azsource.NewMapper(hours)

	return mapper.AsPlanResponse(location, req.Units)
}

func isValidDataRange(dataRange *weatherpb.DataRange, val float64) bool {
	if dataRange == nil {
		return true
	}
	return val >= dataRange.Min && val <= dataRange.Max
}

func (s server) lookupLocation(lookup string) (*geocodingpb.GeocodeReply, error) {
	geocodeRequest := &geocodingpb.GeocodeRequest{
		Location: lookup,
	}

	location, err := s.GeocodeClient.Geocode(context.Background(), geocodeRequest)
	if err != nil {
		return nil, errors.New("geocode client error: " + err.Error())
	}

	return location, nil
}

func unitsAreValid(units string) bool {
	return units == "us" || units == "si"
}
