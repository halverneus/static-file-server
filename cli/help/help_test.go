package help

import "testing"

func TestRun(t *testing.T) {
	if err := Run(); nil != err {
		t.Errorf("While running help got %v", err)
	}
}
