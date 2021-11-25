package itemcondition

type ItemCondition int64

const (
	New ItemCondition = iota
	Used
	PreOwned
	Refurbished
	Unknown
)

func (i ItemCondition) String() string {
	switch i {
	case New:
		return "New"
	case Used:
		return "Used"
	case PreOwned:
		return "PreOwned"
	case Refurbished:
		return "Refurbished"
	}
	return "Unknown"
}

func ParseItemCondition(s string) ItemCondition {
	switch s {
	case "New":
		return New
	case "Used":
		return Used
	case "PreOwned":
		return PreOwned
	case "Refurbished":
		return Refurbished
	}
	return Unknown
}
