package units

type Units interface {
	Distance() string
	Speed() string
	Temperature() string
	Accumulation() string
	Pressure() string
}
