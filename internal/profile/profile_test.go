package profile_test

import (
	"os"
	"path"
	"testing"

	"github.com/quinn-tao/hmis/v1/internal/profile"
)

// Test that encoding a profile should decode to the same profile
func TestProfileEncodeDecode(t *testing.T) {
    testProfilePath := "test_profile.yaml"
    var p profile.Profile
    p.ReadFrom(testProfilePath)
    tmpPath := path.Join(os.TempDir(), "test_profile.yaml")
    p.WriteBack(tmpPath)
    var p1 profile.Profile
    p1.ReadFrom(tmpPath)
    cmpProfile(t, &p, &p1)
}

// Only used for testing purpose
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
