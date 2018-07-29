package cli

import (
	"testing"
)

func TestParse(t *testing.T) {
	matches := func(args Args, orig []string) bool {
		if nil == orig {
			return nil == args
		}
		if len(orig) != len(args) {
			return false
		}
		for index, value := range args {
			if orig[index] != value {
				return false
			}
		}
		return true
	}

	testCases := []struct {
		name  string
		value []string
	}{
		{"Nil arguments", nil},
		{"No arguments", []string{}},
		{"Arguments", []string{"first", "second", "*"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if args := Parse(tc.value); !matches(args, tc.value) {
				t.Errorf("Expected [%v] but got [%v]", tc.value, args)
			}
		})
	}
}

func TestMatches(t *testing.T) {
	testCases := []struct {
		name    string
		value   []string
		pattern []string
		result  bool
	}{
		{"Nil args and nil pattern", nil, nil, true},
		{"No args and nil pattern", []string{}, nil, true},
		{"Nil args and no pattern", nil, []string{}, true},
		{"No args and no pattern", []string{}, []string{}, true},
		{"Nil args and pattern", nil, []string{"test"}, false},
		{"No args and pattern", []string{}, []string{"test"}, false},
		{"Args and nil pattern", []string{"test"}, nil, false},
		{"Args and no pattern", []string{"test"}, []string{}, false},
		{"Simple single compare", []string{"test"}, []string{"test"}, true},
		{"Simple double compare", []string{"one", "two"}, []string{"one", "two"}, true},
		{"Bad single", []string{"one"}, []string{"two"}, false},
		{"Bad double", []string{"one", "two"}, []string{"one", "owt"}, false},
		{"Count mismatch", []string{"one", "two"}, []string{"one"}, false},
		{"Nil args and wild", nil, []string{"*"}, false},
		{"No args and wild", []string{}, []string{"*"}, false},
		{"Single arg and wild", []string{"one"}, []string{"*"}, true},
		{"Double arg and first wild", []string{"one", "two"}, []string{"*", "two"}, true},
		{"Double arg and second wild", []string{"one", "two"}, []string{"one", "*"}, true},
		{"Double arg and first wild mismatched", []string{"one", "two"}, []string{"*", "owt"}, false},
		{"Double arg and second wild mismatched", []string{"one", "two"}, []string{"eno", "*"}, false},
		{"Double arg and double wild", []string{"one", "two"}, []string{"*", "*"}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := Parse(tc.value)
			if resp := args.Matches(tc.pattern...); tc.result != resp {
				msg := "For arguments [%v] matched to pattern [%v] expected " +
					"%b but got %b"
				t.Errorf(msg, tc.value, tc.pattern, tc.result, resp)
			}
		})
	}
}
