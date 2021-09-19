package azsource

import (
	"fmt"
	"strings"

	"github.com/shawntoffel/azure-maps-go/azweather"
)

type Formatter struct {
}

func NewForecastFormatter() Formatter {
	return Formatter{}
}

func (f *Formatter) WeatherUnit(w *azweather.WeatherUnit) string {
	return fmt.Sprintf("%g %s", w.Value, w.Unit)
}

func (f *Formatter) Title(title string) string {
	return strings.Title(title)
}

func (f *Formatter) SpeedWithBearing(w *azweather.Wind) string {
	formattedSpeed := f.WeatherUnit(w.Speed)
	return fmt.Sprintf("%s %s", formattedSpeed, w.Direction.LocalizedDescription)
}
