package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestLoad(t *testing.T) {
	// Verify envvars are set.
	testFolder := "/my/directory"
	os.Setenv(folderKey, testFolder)
	if err := Load(""); nil != err {
		t.Errorf(
			"While loading an empty file name expected no error but got %v",
			err,
		)
	}
	if Get.Folder != testFolder {
		t.Errorf(
			"While loading an empty file name expected folder %s but got %s",
			testFolder, Get.Folder,
		)
	}

	// Verify error if file doesn't exist.
	if err := Load("/this/file/should/never/exist"); nil == err {
		t.Error("While loading non-existing file expected error but got nil")
	}

	// Verify bad YAML returns an error.
	func(t *testing.T) {
		filename := "testing.tmp"
		contents := []byte("{")
		defer os.Remove(filename)

		if err := ioutil.WriteFile(filename, contents, 0666); nil != err {
			t.Errorf("Failed to save bad YAML file with: %v\n", err)
		}
		if err := Load(filename); nil == err {
			t.Error("While loading bad YAML expected error but got nil")
		}
	}(t)

	// Verify good YAML returns no error and sets value.
	func(t *testing.T) {
		filename := "testing.tmp"
		testFolder := "/test/folder"
		contents := []byte(fmt.Sprintf(
			`{"folder": "%s"}`, testFolder,
		))
		defer os.Remove(filename)

		if err := ioutil.WriteFile(filename, contents, 0666); nil != err {
			t.Errorf("Failed to save good YAML file with: %v\n", err)
		}
		if err := Load(filename); nil != err {
			t.Errorf(
				"While loading good YAML expected nil but got %v",
				err,
			)
		}
	}(t)
}

func TestLog(t *testing.T) {
	// Test whether YAML marshalling works, as that is the only error case.
	if _, err := yaml.Marshal(&Get); nil != err {
		t.Errorf("While testing YAML marshalling for config Log() got %v", err)
	}
	Log()
}

func TestOverrideWithEnvvars(t *testing.T) {
	// Choose values that are different than defaults.
	testDebug := true
	testFolder := "/my/directory"
	testHost := "apets.life"
	testPort := uint16(666)
	testShowListing := false
	testTLSCert := "my.pem"
	testTLSKey := "my.key"
	testURLPrefix := "/url/prefix"

	// Set all environment variables with test values.
	os.Setenv(debugKey, fmt.Sprintf("%t", testDebug))
	os.Setenv(folderKey, testFolder)
	os.Setenv(hostKey, testHost)
	os.Setenv(portKey, strconv.Itoa(int(testPort)))
	os.Setenv(showListingKey, fmt.Sprintf("%t", testShowListing))
	os.Setenv(tlsCertKey, testTLSCert)
	os.Setenv(tlsKeyKey, testTLSKey)
	os.Setenv(urlPrefixKey, testURLPrefix)

	// Verification functions.
	equalStrings := func(t *testing.T, name, key, expected, result string) {
		if expected != result {
			t.Errorf(
				"While checking %s for '%s' expected '%s' but got '%s'",
				name, key, expected, result,
			)
		}
	}
	equalUint16 := func(t *testing.T, name, key string, expected, result uint16) {
		if expected != result {
			t.Errorf(
				"While checking %s for '%s' expected %d but got %d",
				name, key, expected, result,
			)
		}
	}
	equalBool := func(t *testing.T, name, key string, expected, result bool) {
		if expected != result {
			t.Errorf(
				"While checking %s for '%s' expected %t but got %t",
				name, key, expected, result,
			)
		}
	}

	// Verify defaults.
	setDefaults()
	phase := "defaults"
	equalBool(t, phase, debugKey, defaultDebug, Get.Debug)
	equalStrings(t, phase, folderKey, defaultFolder, Get.Folder)
	equalStrings(t, phase, hostKey, defaultHost, Get.Host)
	equalUint16(t, phase, portKey, defaultPort, Get.Port)
	equalBool(t, phase, showListingKey, defaultShowListing, Get.ShowListing)
	equalStrings(t, phase, tlsCertKey, defaultTLSCert, Get.TLSCert)
	equalStrings(t, phase, tlsKeyKey, defaultTLSKey, Get.TLSKey)
	equalStrings(t, phase, urlPrefixKey, defaultURLPrefix, Get.URLPrefix)

	// Apply overrides.
	overrideWithEnvVars()

	// Verify overrides.
	phase = "overrides"
	equalBool(t, phase, debugKey, testDebug, Get.Debug)
	equalStrings(t, phase, folderKey, testFolder, Get.Folder)
	equalStrings(t, phase, hostKey, testHost, Get.Host)
	equalUint16(t, phase, portKey, testPort, Get.Port)
	equalBool(t, phase, showListingKey, testShowListing, Get.ShowListing)
	equalStrings(t, phase, tlsCertKey, testTLSCert, Get.TLSCert)
	equalStrings(t, phase, tlsKeyKey, testTLSKey, Get.TLSKey)
	equalStrings(t, phase, urlPrefixKey, testURLPrefix, Get.URLPrefix)
}

func TestValidate(t *testing.T) {
	validPath := "config.go"
	invalidPath := "should/never/exist.txt"
	empty := ""
	prefix := "/my/prefix"

	testCases := []struct {
		name    string
		cert    string
		key     string
		prefix  string
		isError bool
	}{
		{"Valid paths w/prefix", validPath, validPath, prefix, false},
		{"Valid paths wo/prefix", validPath, validPath, empty, false},
		{"Empty paths w/prefix", empty, empty, prefix, false},
		{"Empty paths wo/prefix", empty, empty, empty, false},
		{"Mixed paths w/prefix", empty, validPath, prefix, true},
		{"Alt mixed paths w/prefix", validPath, empty, prefix, true},
		{"Mixed paths wo/prefix", empty, validPath, empty, true},
		{"Alt mixed paths wo/prefix", validPath, empty, empty, true},
		{"Invalid cert w/prefix", invalidPath, validPath, prefix, true},
		{"Invalid key w/prefix", validPath, invalidPath, prefix, true},
		{"Invalid cert & key w/prefix", invalidPath, invalidPath, prefix, true},
		{"Prefix missing leading /", empty, empty, "my/prefix", true},
		{"Prefix with trailing /", empty, empty, "/my/prefix/", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			Get.TLSCert = tc.cert
			Get.TLSKey = tc.key
			Get.URLPrefix = tc.prefix
			err := validate()
			hasError := nil != err
			if hasError && !tc.isError {
				t.Errorf("Expected no error but got %v", err)
			}
			if !hasError && tc.isError {
				t.Error("Expected an error but got no error")
			}
		})
	}
}

func TestEnvAsStr(t *testing.T) {
	sv := "STRING_VALUE"
	fv := "FLOAT_VALUE"
	iv := "INT_VALUE"
	bv := "BOOL_VALUE"
	ev := "EMPTY_VALUE"
	uv := "UNSET_VALUE"

	sr := "String Cheese"    // String result
	fr := "123.456"          // Float result
	ir := "-123"             // Int result
	br := "true"             // Bool result
	er := ""                 // Empty result
	fbr := "fallback result" // Fallback result
	efbr := ""               // Empty fallback result

	os.Setenv(sv, sr)
	os.Setenv(fv, fr)
	os.Setenv(iv, ir)
	os.Setenv(bv, br)
	os.Setenv(ev, er)

	testCases := []struct {
		name     string
		key      string
		fallback string
		result   string
	}{
		{"Good string", sv, fbr, sr},
		{"Float string", fv, fbr, fr},
		{"Int string", iv, fbr, ir},
		{"Bool string", bv, fbr, br},
		{"Empty string", ev, fbr, fbr},
		{"Unset", uv, fbr, fbr},
		{"Good string with empty fallback", sv, efbr, sr},
		{"Unset with empty fallback", uv, efbr, efbr},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := envAsStr(tc.key, tc.fallback)
			if tc.result != result {
				t.Errorf(
					"For %s with a '%s' fallback expected '%s' but got '%s'",
					tc.key, tc.fallback, tc.result, result,
				)
			}
		})
	}
}

func TestEnvAsStrSlice(t *testing.T) {
	oe := "ONE_ENTRY"
	oewc := "ONE_ENTRY_WITH_COMMA"
	oewtc := "ONE_ENTRY_WITH_TRAILING_COMMA"
	te := "TWO_ENTRY"
	tewc := "TWO_ENTRY_WITH_COMMA"
	oc := "ONLY_COMMA"
	ev := "EMPTY_VALUE"
	uv := "UNSET_VALUE"

	fs := "http://my.site"
	ts := "http://other.site"
	fbr := []string{"one", "two"}
	var efbr []string

	oes := fs
	oer := []string{fs}
	oewcs := "," + fs
	oewcr := []string{"", fs}
	oewtcs := fs + ","
	oewtcr := []string{fs, ""}
	tes := fs + "," + ts
	ter := []string{fs, ts}
	tewcs := "," + fs + "," + ts
	tewcr := []string{"", fs, ts}
	ocs := ","
	ocr := []string{"", ""}
	evs := ""

	os.Setenv(oe, oes)
	os.Setenv(oewc, oewcs)
	os.Setenv(oewtc, oewtcs)
	os.Setenv(te, tes)
	os.Setenv(tewc, tewcs)
	os.Setenv(oc, ocs)
	os.Setenv(ev, evs)

	testCases := []struct {
		name     string
		key      string
		fallback []string
		result   []string
	}{
		{"One entry", oe, fbr, oer},
		{"One entry w/comma", oewc, fbr, oewcr},
		{"One entry w/trailing comma", oewtc, fbr, oewtcr},
		{"Two entry", te, fbr, ter},
		{"Two entry w/comma", tewc, fbr, tewcr},
		{"Only comma", oc, fbr, ocr},
		{"Empty value w/fallback", ev, fbr, fbr},
		{"Empty value wo/fallback", ev, efbr, efbr},
		{"Unset w/fallback", uv, fbr, fbr},
		{"Unset wo/fallback", uv, efbr, efbr},
	}

	matches := func(a, b []string) bool {
		if len(a) != len(b) {
			return false
		}
		tally := make(map[int]bool)
		for i := range a {
			tally[i] = false
		}
		for _, val := range a {
			for i, other := range b {
				if other == val && !tally[i] {
					tally[i] = true
					break
				}
			}
		}
		for _, found := range tally {
			if !found {
				return false
			}
		}
		return true
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := envAsStrSlice(tc.key, tc.fallback)
			if !matches(tc.result, result) {
				t.Errorf(
					"For %s with a '%v' fallback expected '%v' but got '%v'",
					tc.key, tc.fallback, tc.result, result,
				)
			}
		})
	}
}

func TestEnvAsUint16(t *testing.T) {
	ubv := "UPPER_BOUNDS_VALUE"
	lbv := "LOWER_BOUNDS_VALUE"
	hv := "HIGH_VALUE"
	lv := "LOW_VALUE"
	bv := "BOOL_VALUE"
	sv := "STRING_VALUE"
	uv := "UNSET_VALUE"

	fbr := uint16(666)   // Fallback result
	ubr := uint16(65535) // Upper bounds result
	lbr := uint16(0)     // Lower bounds result

	os.Setenv(ubv, "65535")
	os.Setenv(lbv, "0")
	os.Setenv(hv, "65536")
	os.Setenv(lv, "-1")
	os.Setenv(bv, "true")
	os.Setenv(sv, "Cheese")

	testCases := []struct {
		name     string
		key      string
		fallback uint16
		result   uint16
	}{
		{"Upper bounds", ubv, fbr, ubr},
		{"Lower bounds", lbv, fbr, lbr},
		{"Out-of-bounds high", hv, fbr, fbr},
		{"Out-of-bounds low", lv, fbr, fbr},
		{"Boolean", bv, fbr, fbr},
		{"String", sv, fbr, fbr},
		{"Unset", uv, fbr, fbr},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := envAsUint16(tc.key, tc.fallback)
			if tc.result != result {
				t.Errorf(
					"For %s with a %d fallback expected %d but got %d",
					tc.key, tc.fallback, tc.result, result,
				)
			}
		})
	}
}

func TestEnvAsBool(t *testing.T) {
	tv := "TRUE_VALUE"
	fv := "FALSE_VALUE"
	bv := "BAD_VALUE"
	uv := "UNSET_VALUE"

	os.Setenv(tv, "True")
	os.Setenv(fv, "NO")
	os.Setenv(bv, "BAD")

	testCases := []struct {
		name     string
		key      string
		fallback bool
		result   bool
	}{
		{"True with true fallback", tv, true, true},
		{"True with false fallback", tv, false, true},
		{"False with true fallback", fv, true, false},
		{"False with false fallback", fv, false, false},
		{"Bad with true fallback", bv, true, true},
		{"Bad with false fallback", bv, false, false},
		{"Unset with true fallback", uv, true, true},
		{"Unset with false fallback", uv, false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := envAsBool(tc.key, tc.fallback)
			if tc.result != result {
				t.Errorf(
					"For %s with a %t fallback expected %t but got %t",
					tc.key, tc.fallback, tc.result, result,
				)
			}
		})
	}
}

func TestStrAsBool(t *testing.T) {
	testCases := []struct {
		name    string
		value   string
		result  bool
		isError bool
	}{
		{"Empty value", "", false, true},
		{"False value", "0", false, false},
		{"True value", "1", true, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := strAsBool(tc.value)
			if result != tc.result {
				t.Errorf(
					"Expected %t for %s but got %t",
					tc.result, tc.value, result,
				)
			}
			if tc.isError && nil == err {
				t.Errorf(
					"Expected error for %s but got no error",
					tc.value,
				)
			}
			if !tc.isError && nil != err {
				t.Errorf(
					"Expected no error for %s but got %v",
					tc.value, err,
				)
			}
		})

	}
}
