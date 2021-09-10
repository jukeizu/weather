package azmapssource

import (
	"fmt"

	gpb "github.com/jukeizu/weather/api/protobuf-spec/geocodingpb"
	wpb "github.com/jukeizu/weather/api/protobuf-spec/weatherpb"
	"github.com/shawntoffel/azure-maps-go/weather/entities"
	"github.com/urkk/metar"
)

type Mapper struct {
	formatter ForecastFormatter
	currently entities.CurrentConditions
}

func NewMapper(currently entities.CurrentConditions) Mapper {
	formatter := NewForecastFormatter()

	return Mapper{formatter, currently}
}

func (m *Mapper) AsWeatherResponse(location *gpb.GeocodeReply) (*wpb.WeatherReply, error) {
	response := &wpb.WeatherReply{}

	response.Location = location.FormattedAddress
	response.Latitude = location.Latitude
	response.Longitude = location.Longitude
	response.Currently = m.mapCurrently()
	//response.Forecast = m.mapForecast()
	//response.Alerts = m.mapAlerts()

	return response, nil
}

func (m *Mapper) mapCurrently() *wpb.Currently {
	currently := &wpb.Currently{}

	dataPoint := m.currently

	currently.Description = dataPoint.Phrase
	//currently.Summary = m.formatter.CombinedDataBlockSummary(m.forecast.Minutely, m.forecast.Hourly)
	currently.Data = m.mapData(dataPoint)

	msg, _ := metar.NewMETAR("KFLY 070415Z AUTO 29011KT 10SM CLR 17/00 A3028 RMK AO2 T01700004")
	currently.Summary = fmt.Sprintf("%d", msg.Wind.SpeedMps())

	return currently
}

func (m *Mapper) mapData(dataPoint entities.CurrentConditions) *wpb.Data {
	data := &wpb.Data{}

	data.DewPoint = m.formatter.WeatherUnit(dataPoint.DewPoint)
	data.FeelsLike = m.formatter.WeatherUnit(dataPoint.RealFeelTemperature)
	//data.Humidity = m.formatter.Percentage(dataPoint.Humidity)
	//data.Icon =
	//data.PrecipitationProbability = m.formatter.Percentage(dataPoint.PrecipProbability)
	//data.PrecipitationType = m.formatter.Title(dataPoint.PrecipType)
	data.Pressure = m.formatter.WeatherUnit(dataPoint.Pressure)
	data.Temperature = m.formatter.WeatherUnit(dataPoint.Temperature)
	//data.TemperatureHigh = m.formatter.Temperature(dataPoint.TemperatureHigh)
	//data.TemperatureLow = m.formatter.Temperature(dataPoint.TemperatureLow)
	//data.Timestamp = dataPoint.DateTime
	data.Wind = m.formatter.SpeedWithBearing(dataPoint.Wind) + " gusting " + m.formatter.WeatherUnit(dataPoint.WindGust.Speed)

	return data
}
