package version

import "testing"

func TestVersion(t *testing.T) {
	if err := Run(); nil != err {
		t.Errorf("While running version got %v", err)
	}
}
