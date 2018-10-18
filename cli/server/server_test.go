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

	testCases := []struct {
		name    string
		folder  string
		prefix  string
		listing bool
		debug   bool
	}{
		{"Basic handler w/o debug", testFolder, "", true, false},
		{"Prefix handler w/o debug", testFolder, testPrefix, true, false},
		{"Basic and hide listing handler w/o debug", testFolder, "", false, false},
		{"Prefix and hide listing handler w/o debug", testFolder, testPrefix, false, false},
		{"Basic handler w/debug", testFolder, "", true, true},
		{"Prefix handler w/debug", testFolder, testPrefix, true, true},
		{"Basic and hide listing handler w/debug", testFolder, "", false, true},
		{"Prefix and hide listing handler w/debug", testFolder, testPrefix, false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config.Get.Debug = tc.debug
			config.Get.Folder = tc.folder
			config.Get.ShowListing = tc.listing
			config.Get.URLPrefix = tc.prefix

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
