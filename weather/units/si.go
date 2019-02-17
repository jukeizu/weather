package units

type SiUnits struct {
	Units
}

func (u SiUnits) Distance() string {
	return "km"
}

func (u SiUnits) Speed() string {
	return "km/h"
}

func (u SiUnits) Temperature() string {
	return "°C"
}

func (u SiUnits) Accumulation() string {
	return "cm"
}

func (u SiUnits) Pressure() string {
	return "hPa"
}
