package units

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
