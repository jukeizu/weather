package darkskysource

type Units interface {
	Distance() string
	Speed() string
	Temperature() string
	Accumulation() string
	Pressure() string
}

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

type UsUnits struct {
	Units
}

func (u UsUnits) Distance() string {
	return "mi"
}

func (u UsUnits) Speed() string {
	return "mph"
}

func (u UsUnits) Temperature() string {
	return "°F"
}

func (u UsUnits) Accumulation() string {
	return "in"
}

func (u UsUnits) Pressure() string {
	return "mb"
}
