package profile_test

import (
	"os"
	"testing"

	"github.com/quinn-tao/hmis/v1/internal/profile"
)

// Test that encoding a profile should decode to the same profile
func TestProfileEncodeDecode(t *testing.T) {
	testProfilePath := "test_profile.yaml"

	var p profile.Profile
	testInputFile, err := os.Open(testProfilePath)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer testInputFile.Close()

	err = p.ReadFrom(testInputFile)
	t.Logf("%v", p)
	if err != nil {
		t.Fatalf("Error while reading profile: %v", err)
	}

	file, err := os.CreateTemp("", "test_profile.yaml")
	err = p.WriteBack(file)
	if err != nil {
		t.Fatalf("Error while writing profile: %v", err)
	}
	file.Close()

	var p1 profile.Profile
	file, err = os.Open(file.Name())
	err = p1.ReadFrom(file)
	if err != nil {
		t.Fatalf("Error while reading from tmp profile: %v", err)
	}
	cmpProfile(t, &p, &p1)
}

func cmpProfile(t *testing.T, exp *profile.Profile, act *profile.Profile) {
	if exp.Name != act.Name {
		t.Fatalf("Expected name %v, got name %v", exp.Name, act.Name)
	}
	if exp.Mode != act.Mode {
		t.Fatalf("Expected mode %v, got mode %v", exp.Mode, act.Mode)
	}
	if !exp.Category.Equals(act.Category) {
		t.Fatalf("Expected categrories %v, got categories %v", exp.Category, act.Category)
	}
}
