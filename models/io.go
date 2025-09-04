package models

type IO int

const (
	Straight IO = iota + 1
	Double
	Master
)

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
