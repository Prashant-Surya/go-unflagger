package flagger

const (
	DATE = "date"
	NAME = "name"
)

type FlagChecker interface {
	IsValidFlag(conditions []string) bool
}

func NewFlagger(flaggerType string, name string, dateFormat string) *CommonFlagger {
	obj := &CommonFlagger{}
	switch flaggerType {
	case DATE:
		obj.FlagCheckerObj = &DateFlagChecker{
			DateFormat: dateFormat,
		}
	case NAME:
		obj.FlagCheckerObj = &NameFlagChecker{
			Name: name,
		}
	default:
		return nil
	}
	return obj
}