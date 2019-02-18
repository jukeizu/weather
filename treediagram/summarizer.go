package treediagram

import (
	"bytes"
	"fmt"

	"github.com/jukeizu/weather/api/protobuf-spec/weatherpb"
)

var iconMap = map[string]string{
	"clear-day":           ":sunny:",
	"clear-night":         ":crescent_moon:",
	"rain":                ":cloud_rain:",
	"snow":                ":cloud_snow:",
	"sleet":               ":cloud_rain:",
	"wind":                ":dash:",
	"fog":                 ":fog:",
	"cloudy":              ":cloud:",
	"partly-cloudy-day":   ":white_sun_cloud:",
	"partly-cloudy-night": ":cloud:",
	"thunderstorm":        ":thunder_cloud_rain:",
	"hail":                ":cloud_rain:",
	"tornado":             ":cloud_tornado:",
}

func generateProbabilitySummary(data *weatherpb.Data) string {
	buffer := bytes.Buffer{}

	precipitationTitle := "Precipitation"

	if data.PrecipitationType != "" {
		precipitationTitle = data.PrecipitationType
	}

	buffer.WriteString(precipitationTitle + ": " + data.PrecipitationProbability)
	buffer.WriteString("\nFeels Like: " + data.FeelsLike)

	return buffer.String()
}

func generateDataSummary(data *weatherpb.Data) string {
	buffer := bytes.Buffer{}

	buffer.WriteString("Humidity: " + data.Humidity)
	buffer.WriteString("\nWind: " + data.Wind)
	buffer.WriteString("\nBarometer: " + data.Pressure)
	buffer.WriteString("\nDew Point: " + data.DewPoint)

	return buffer.String()
}

func generateDayTitle(day *weatherpb.Day) string {
	buffer := bytes.Buffer{}

	data := day.Data

	buffer.WriteString(getEmojiForIcon(data.Icon))
	buffer.WriteString(day.Weekday)

	if data.TemperatureHigh != "" {
		buffer.WriteString(fmt.Sprintf(" (%s)", data.TemperatureHigh))
	}

	return buffer.String()
}

func getEmojiForIcon(icon string) string {
	emoji, found := iconMap[icon]
	if !found {
		return ""
	}

	return emoji + " "
}
