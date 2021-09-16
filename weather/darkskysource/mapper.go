package darkskysource

import (
	gpb "github.com/jukeizu/weather/api/protobuf-spec/geocodingpb"
	wpb "github.com/jukeizu/weather/api/protobuf-spec/weatherpb"
	"github.com/shawntoffel/darksky"
)

type Mapper struct {
	formatter ForecastFormatter
	forecast  darksky.ForecastResponse
}

func NewMapper(forecast darksky.ForecastResponse) Mapper {
	formatter := NewForecastFormatter(forecast.Flags)

	return Mapper{formatter, forecast}
}

func (m *Mapper) AsWeatherResponse(location *gpb.GeocodeReply) (*wpb.WeatherReply, error) {
	response := &wpb.WeatherReply{}

	response.Location = location.FormattedAddress
	response.Latitude = m.formatter.Coordinate(m.forecast.Latitude)
	response.Longitude = m.formatter.Coordinate(m.forecast.Longitude)
	response.Currently = m.mapCurrently()
	response.Forecast = m.mapForecast()
	response.Alerts = m.mapAlerts()

	return response, nil
}

func (m *Mapper) mapCurrently() *wpb.Currently {
	currently := &wpb.Currently{}

	dataPoint := m.forecast.Currently

	if dataPoint == nil {
		return currently
	}

	currently.Description = dataPoint.Summary
	currently.Summary = m.formatter.CombinedDataBlockSummary(m.forecast.Minutely, m.forecast.Hourly)
	currently.Data = m.mapData(dataPoint)

	return currently
}

func (m *Mapper) mapForecast() *wpb.Forecast {
	forecast := &wpb.Forecast{}

	forecast.Summary = m.formatter.DataBlockSummary(m.forecast.Daily)
	forecast.Days = m.mapDays()

	return forecast
}

func (m *Mapper) mapDays() []*wpb.Day {
	days := []*wpb.Day{}

	dailyForecast := m.forecast.Daily

	if dailyForecast == nil || len(dailyForecast.Data) < 1 {
		return days
	}

	for _, forecastDay := range dailyForecast.Data {
		days = append(days, m.mapDay(forecastDay))
	}

	return days
}

func (m *Mapper) mapDay(dataPoint darksky.DataPoint) *wpb.Day {
	day := &wpb.Day{}

	day.Weekday = m.formatter.Day(dataPoint.Time, m.forecast.Timezone)
	day.Summary = dataPoint.Summary
	day.Data = m.mapData(&dataPoint)

	return day
}

func (m *Mapper) mapAlerts() []*wpb.Alert {
	alerts := []*wpb.Alert{}

	forecastAlerts := m.forecast.Alerts

	if len(forecastAlerts) < 1 {
		return alerts
	}

	for _, forecastAlert := range forecastAlerts {
		alerts = append(alerts, m.mapAlert(forecastAlert))
	}

	return alerts
}

func (m *Mapper) mapAlert(forecastAlert *darksky.Alert) *wpb.Alert {
	alert := &wpb.Alert{}

	if forecastAlert == nil {
		return alert
	}

	alert.Message = m.formatter.AlertMessage(forecastAlert, m.forecast.Timezone)
	alert.Severity = forecastAlert.Severity
	alert.Uri = forecastAlert.Uri

	return alert
}

func (m *Mapper) mapData(dataPoint *darksky.DataPoint) *wpb.Data {
	data := &wpb.Data{}

	if dataPoint == nil {
		return data
	}

	data.DewPoint = m.formatter.Temperature(dataPoint.DewPoint)
	data.FeelsLike = m.formatter.Temperature(dataPoint.ApparentTemperature)
	data.Humidity = m.formatter.Percentage(dataPoint.Humidity)
	data.Icon = dataPoint.Icon
	data.PrecipitationProbability = m.formatter.Percentage(dataPoint.PrecipProbability)
	data.PrecipitationType = m.formatter.Title(dataPoint.PrecipType)
	data.Pressure = m.formatter.Pressure(dataPoint.Pressure)
	data.Temperature = m.formatter.Temperature(dataPoint.Temperature)
	data.TemperatureHigh = m.formatter.Temperature(dataPoint.TemperatureHigh)
	data.TemperatureLow = m.formatter.Temperature(dataPoint.TemperatureLow)
	data.Timestamp = int64(dataPoint.Time)
	data.Wind = m.formatter.SpeedWithBearing(dataPoint.WindSpeed, dataPoint.WindBearing)

	return data
}
