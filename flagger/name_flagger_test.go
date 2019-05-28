package flagger

import "testing"

func TestIsValidFlag_Valid(t *testing.T) {
	var flagChecker = &NameFlagChecker{
		Name: "Test",
	}
	conditions := []string{"Test", "FeatureFlags"}
	if !flagChecker.IsValidFlag(conditions) {
		t.Errorf("Flag checker returned false instead of true")
	}
}

func TestIsValidFlag_InValid(t *testing.T) {
	var flagChecker = &NameFlagChecker{
		Name: "Test",
	}
	conditions := []string{"Tes", "FeatureFlags"}
	if flagChecker.IsValidFlag(conditions) {
		t.Errorf("Flag checker returned true instead of false")
	}
}