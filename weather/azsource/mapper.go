package azsource

import (
	"strconv"
	"time"

	gpb "github.com/jukeizu/weather/api/protobuf-spec/geocodingpb"
	wpb "github.com/jukeizu/weather/api/protobuf-spec/weatherpb"
	"github.com/shawntoffel/azure-maps-go/azweather"
)

type Mapper struct {
	formatter Formatter
	hourly    []azweather.HourlyForecast
}

func NewMapper(hourly []azweather.HourlyForecast) Mapper {
	formatter := NewForecastFormatter()

	return Mapper{formatter, hourly}
}

func (m *Mapper) AsPlanResponse(location *gpb.GeocodeReply, units string) (*wpb.PlanReply, error) {
	response := &wpb.PlanReply{}

	response.Location = location.FormattedAddress
	response.Latitude = location.Latitude
	response.Longitude = location.Longitude
	response.Hours = m.mapHours()
	response.Units = units
	response.GeneratedAt = time.Now().UTC().Unix()

	return response, nil
}

func (m *Mapper) mapHours() []*wpb.Hour {
	hours := []*wpb.Hour{}

	for _, hour := range m.hourly {
		hours = append(hours, m.mapHour(hour))
	}

	return hours
}

func (m *Mapper) mapHour(dataPoint azweather.HourlyForecast) *wpb.Hour {
	return &wpb.Hour{
		Data: m.mapData(dataPoint),
	}
}

func (m *Mapper) mapData(dataPoint azweather.HourlyForecast) *wpb.Data {
	data := &wpb.Data{}
	data.Timestamp = m.mapTime(dataPoint.Date).Unix()
	data.DewPoint = m.formatter.WeatherUnit(dataPoint.DewPoint)
	data.FeelsLike = m.formatter.WeatherUnit(dataPoint.RealFeelTemperature)
	data.Humidity = strconv.Itoa(dataPoint.RelativeHumidity)
	data.PrecipitationProbability = strconv.Itoa(dataPoint.PrecipitationProbability)
	data.Temperature = m.formatter.WeatherUnit(dataPoint.Temperature)
	data.Wind = m.formatter.WeatherUnit(dataPoint.Wind.Speed) + " gusting " + m.formatter.WeatherUnit(dataPoint.WindGust.Speed)

	return data
}

func (m *Mapper) mapTime(val string) time.Time {
	t, _ := time.Parse(time.RFC3339, val)
	return t
}
