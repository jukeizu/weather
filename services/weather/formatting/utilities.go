package formatting

import (
	"github.com/jukeizu/weather/services/weather/units"
	"github.com/shawntoffel/darksky"
)

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
