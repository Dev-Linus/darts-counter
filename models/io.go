package models

// IO represents the in/out mode for starting or finishing a match (Straight, Double, Master).
type IO int

const (
	// Straight allows any throw for in/out.
	Straight IO = iota + 1
	// Double requires doubles for in/out.
	Double
	// Master requires doubles or triples for in/out.
	Master
)

// MapNumberToIO maps a numeric code to its corresponding IO enum value.
func MapNumberToIO(integer uint8) IO {
	switch integer {
	case 0:
		return Straight
	case 1:
		return Double
	case 2:
		return Master
	default:
		return Straight
	}
}

// MapIOToNumber maps an IO enum value to its numeric code.
func MapIOToNumber(io IO) uint8 {
	switch io {
	case Straight:
		return 0
	case Double:
		return 1
	case Master:
		return 2
	default:
		return 0
	}
}

// GetAllFinishingThrows returns the allowed last-throw types for the given IO mode.
func (io IO) GetAllFinishingThrows() []ThrowType {
	switch io {
	case Straight:
		return GetAllThrowTypes(true, false, false)
	case Double:
		return GetAllThrowTypes(false, true, false)
	case Master:
		return GetAllThrowTypes(false, false, true)
	default:
		return nil
	}
}
