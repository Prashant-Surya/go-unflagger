package flagger

type NameFlagChecker struct {
	Name string
}

func(NameFlagChecker *NameFlagChecker) IsValidFlag(conditions []string) bool {
	for i, v := range conditions {
		if v == "FeatureFlags" && conditions[i-1] == NameFlagChecker.Name {
			return true
		}
	}
	return false
}