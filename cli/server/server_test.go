package server

import (
	"errors"
	"net/http"
	"testing"

	"github.com/halverneus/static-file-server/config"
	"github.com/halverneus/static-file-server/handle"
)

func TestRun(t *testing.T) {
	listenerError := errors.New("listener")
	selectListener = func() handle.ListenerFunc {
		return func(string, http.HandlerFunc) error {
			return listenerError
		}
	}

	config.Get.Debug = false
	if err := Run(); listenerError != err {
		t.Errorf("Without debug expected %v but got %v", listenerError, err)
	}

	config.Get.Debug = true
	if err := Run(); listenerError != err {
		t.Errorf("With debug expected %v but got %v", listenerError, err)
	}
}

func TestHandlerSelector(t *testing.T) {
	// This test only exercises function branches.
	testFolder := "/web"
	testPrefix := "/url/prefix"
	var ignoreReferrer []string
	testReferrer := []string{"http://localhost"}

	testCases := []struct {
		name    string
		folder  string
		prefix  string
		listing bool
		debug   bool
		refer   []string
	}{
		{"Basic handler w/o debug", testFolder, "", true, false, ignoreReferrer},
		{"Prefix handler w/o debug", testFolder, testPrefix, true, false, ignoreReferrer},
		{"Basic and hide listing handler w/o debug", testFolder, "", false, false, ignoreReferrer},
		{"Prefix and hide listing handler w/o debug", testFolder, testPrefix, false, false, ignoreReferrer},
		{"Basic handler w/debug", testFolder, "", true, true, ignoreReferrer},
		{"Prefix handler w/debug", testFolder, testPrefix, true, true, ignoreReferrer},
		{"Basic and hide listing handler w/debug", testFolder, "", false, true, ignoreReferrer},
		{"Prefix and hide listing handler w/debug", testFolder, testPrefix, false, true, ignoreReferrer},
		{"Basic handler w/o debug w/refer", testFolder, "", true, false, testReferrer},
		{"Prefix handler w/o debug w/refer", testFolder, testPrefix, true, false, testReferrer},
		{"Basic and hide listing handler w/o debug w/refer", testFolder, "", false, false, testReferrer},
		{"Prefix and hide listing handler w/o debug w/refer", testFolder, testPrefix, false, false, testReferrer},
		{"Basic handler w/debug w/refer", testFolder, "", true, true, testReferrer},
		{"Prefix handler w/debug w/refer", testFolder, testPrefix, true, true, testReferrer},
		{"Basic and hide listing handler w/debug w/refer", testFolder, "", false, true, testReferrer},
		{"Prefix and hide listing handler w/debug w/refer", testFolder, testPrefix, false, true, testReferrer},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config.Get.Debug = tc.debug
			config.Get.Folder = tc.folder
			config.Get.ShowListing = tc.listing
			config.Get.URLPrefix = tc.prefix
			config.Get.Referrers = tc.refer

			handlerSelector()
		})
	}
}

func TestListenerSelector(t *testing.T) {
	// This test only exercises function branches.
	testCert := "file.crt"
	testKey := "file.key"

	testCases := []struct {
		name string
		cert string
		key  string
	}{
		{"HTTP", "", ""},
		{"HTTPS", testCert, testKey},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config.Get.TLSCert = tc.cert
			config.Get.TLSKey = tc.key
			listenerSelector()
		})
	}
}
