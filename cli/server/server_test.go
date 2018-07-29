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

	if err := Run(); listenerError != err {
		t.Errorf("Expected %v but got %v", listenerError, err)
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
	}{
		{"Basic handler", testFolder, "", true},
		{"Prefix handler", testFolder, testPrefix, true},
		{"Basic and hide listing handler", testFolder, "", false},
		{"Prefix and hide listing handler", testFolder, testPrefix, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config.Get.Folder = tc.folder
			config.Get.URLPrefix = tc.prefix
			config.Get.ShowListing = tc.listing

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
