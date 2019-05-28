package flagger

import "testing"

func TestDateCheckerIsValidFlag_Valid(t *testing.T) {
	var flagChecker = &DateFlagChecker{
		DateFormat: "2006_01_02",
	}
	conditions := []string{"Remove__2018_01_01", "FeatureFlags"}
	if !flagChecker.IsValidFlag(conditions) {
		t.Errorf("Flag checker returned false instead of true")
	}
}

func TestDateCheckerIsValidFlag_InValid(t *testing.T) {
	var flagChecker = &DateFlagChecker{
		DateFormat: "2008_01_02",
	}
	conditions := []string{"Remove__2020_01_01", "FeatureFlags"}
	if flagChecker.IsValidFlag(conditions) {
		t.Errorf("Flag checker returned true instead of false")
	}
}

func TestDateCheckerIsValidFlag_NoDate(t *testing.T) {
	var flagChecker = &DateFlagChecker{
		DateFormat: "2008_01_02",
	}
	conditions := []string{"Remove", "FeatureFlags"}
	if flagChecker.IsValidFlag(conditions) {
		t.Errorf("Flag checker returned true instead of false")
	}
}
