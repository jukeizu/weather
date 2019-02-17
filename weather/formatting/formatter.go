package formatting

import (
	"fmt"
	"strings"
	"time"

	"github.com/jukeizu/weather/weather/units"
	"github.com/shawntoffel/darksky"
)

type ForecastFormatter struct {
	Units units.Units
}

func NewForecastFormatter(flags *darksky.Flags) ForecastFormatter {
	u := getUnits(flags)

	return ForecastFormatter{u}
}

func (f *ForecastFormatter) DataBlockSummary(dataBlock *darksky.DataBlock) string {
	if dataBlock == nil {
		return ""
	}

	return dataBlock.Summary
}

func (f *ForecastFormatter) CombinedDataBlockSummary(first *darksky.DataBlock, second *darksky.DataBlock) string {
	firstSummary := f.DataBlockSummary(first)
	secondSummary := f.DataBlockSummary(second)

	if firstSummary == "" {
		return secondSummary
	}

	if secondSummary == "" {
		return firstSummary
	}

	return firstSummary + " " + secondSummary
}

func (f *ForecastFormatter) Temperature(m darksky.Measurement) string {
	if m == 0 {
		return ""
	}
	return f.Measurement(m) + f.Units.Temperature()
}

func (f *ForecastFormatter) Measurement(m darksky.Measurement) string {
	return fmt.Sprintf("%.1f", m)
}

func (f *ForecastFormatter) Percentage(m darksky.Measurement) string {
	return fmt.Sprintf("%.0f%%", float64(m)*100)
}

func (f *ForecastFormatter) Coordinate(m darksky.Measurement) float64 {
	return float64(m)
}

func (f *ForecastFormatter) Distance(m darksky.Measurement) string {
	return fmt.Sprintf("%s %s", f.Measurement(m), f.Units.Distance())
}

func (f *ForecastFormatter) Pressure(m darksky.Measurement) string {
	return fmt.Sprintf("%s %s", f.Measurement(m), f.Units.Pressure())
}

func (f *ForecastFormatter) Speed(m darksky.Measurement) string {
	return fmt.Sprintf("%s %s", f.Measurement(m), f.Units.Speed())
}

func (f *ForecastFormatter) Time(val darksky.Timestamp, timezone string) string {
	loc, _ := time.LoadLocation(timezone)

	return time.Unix(int64(val), 0).In(loc).Format("3:04 PM MST")
}

func (f *ForecastFormatter) Day(val darksky.Timestamp, timezone string) string {
	loc, _ := time.LoadLocation(timezone)

	return time.Unix(int64(val), 0).In(loc).Format("Mon")
}

func (f *ForecastFormatter) Title(title string) string {
	return strings.Title(title)
}

func (f *ForecastFormatter) AlertMessage(alert *darksky.Alert, timezone string) string {
	if alert == nil {
		return ""
	}

	return fmt.Sprintf("%s until %s", alert.Title, f.Time(alert.Expires, timezone))
}

func (f *ForecastFormatter) SpeedWithBearing(speed darksky.Measurement, bearing darksky.Measurement) string {
	formattedSpeed := f.Speed(speed)

	if bearing == 0 {
		return formattedSpeed
	}

	return fmt.Sprintf("%s %s", f.Bearing(bearing), formattedSpeed)
}

func (f *ForecastFormatter) Bearing(m darksky.Measurement) string {
	switch {
	case m >= 0 && m < 11.25:
		return "N"
	case m >= 11.25 && m < 33.75:
		return "NNE"
	case m >= 33.75 && m < 56.25:
		return "NE"
	case m >= 56.25 && m < 78.75:
		return "ENE"
	case m >= 78.75 && m < 101.25:
		return "E"
	case m >= 101.25 && m < 123.75:
		return "ESE"
	case m >= 123.75 && m < 146.25:
		return "SE"
	case m >= 146.25 && m < 168.75:
		return "SSE"
	case m >= 168.75 && m < 191.25:
		return "S"
	case m >= 191.25 && m < 213.75:
		return "SSW"
	case m >= 213.75 && m < 236.25:
		return "SW"
	case m >= 236.25 && m < 258.75:
		return "WSW"
	case m >= 258.75 && m < 281.25:
		return "W"
	case m >= 281.25 && m < 303.75:
		return "WNW"
	case m >= 303.75 && m < 326.25:
		return "NW"
	case m >= 326.25 && m < 348.75:
		return "NNW"
	case m >= 348.75 && m <= 360.00:
		return "N"
	}

	return ""
}

func getUnits(flags *darksky.Flags) units.Units {

	if flags == nil {
		return units.UsUnits{}
	}

	switch flags.Units {
	case "us":
		return units.UsUnits{}
	case "si":
		return units.SiUnits{}
	default:
		return units.UsUnits{}
	}
}
