package flagger

import (
	"strings"
	"time"
)

type DateFlagChecker struct {
	DateFormat string
}

func(dateFlagger *DateFlagChecker) IsValidFlag(conditions []string) bool {
	var flag string

	for i, v := range conditions {
		if v == "FeatureFlags" {
			flag = conditions[i-1]
			break
		}
	}

	flagsSplit := strings.Split(flag, "__")
	if len(flagsSplit) <= 1 {
		return false
	}

	date := extractDate(dateFlagger.DateFormat, flagsSplit[len(flagsSplit) - 1])
	if date != nil {
		now := time.Now()
		if now.After(*date) {
			return true
		}
	}
	return false
}