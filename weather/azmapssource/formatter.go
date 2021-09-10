package azmapssource

import (
	"fmt"
	"strings"
	"time"

	"github.com/shawntoffel/azure-maps-go/weather/entities"
	"github.com/shawntoffel/darksky"
)

type ForecastFormatter struct {
}

func NewForecastFormatter() ForecastFormatter {
	return ForecastFormatter{}
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

func (f *ForecastFormatter) WeatherUnit(w *entities.WeatherUnit) string {
	return fmt.Sprintf("%g %s", w.Value, w.Unit)
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

func (f *ForecastFormatter) SpeedWithBearing(w *entities.Wind) string {
	formattedSpeed := f.WeatherUnit(w.Speed)
	return fmt.Sprintf("%s %s", w.Direction.LocalizedDescription, formattedSpeed)
}
