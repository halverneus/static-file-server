package server

import (
	"fmt"
	"net/http"

	"github.com/halverneus/static-file-server/config"
	"github.com/halverneus/static-file-server/handle"
)

var (
	// Values to be overridden to simplify unit testing.
	selectHandler  = handlerSelector
	selectListener = listenerSelector
)

// Run server.
func Run() error {
	if config.Get.Debug {
		config.Log()
	}
	// Choose and set the appropriate, optimized static file serving function.
	handler := selectHandler()

	// Serve files over HTTP or HTTPS based on paths to TLS files being
	// provided.
	listener := selectListener()

	binding := fmt.Sprintf("%s:%d", config.Get.Host, config.Get.Port)
	return listener(binding, handler)
}

// handlerSelector returns the appropriate request handler based on
// configuration.
func handlerSelector() (handler http.HandlerFunc) {
	var serveFileHandler handle.FileServerFunc

	serveFileHandler = http.ServeFile
	if config.Get.Debug {
		serveFileHandler = handle.WithLogging(serveFileHandler)
	}

	if 0 != len(config.Get.Referrers) {
		serveFileHandler = handle.WithReferrers(
			serveFileHandler, config.Get.Referrers,
		)
	}

	// Choose and set the appropriate, optimized static file serving function.
	if 0 == len(config.Get.URLPrefix) {
		handler = handle.Basic(serveFileHandler, config.Get.Folder)
	} else {
		handler = handle.Prefix(
			serveFileHandler,
			config.Get.Folder,
			config.Get.URLPrefix,
		)
	}

	// Determine whether index files should hidden.
	if !config.Get.ShowListing {
		if config.Get.AllowIndex {
			handler = handle.PreventListings(handler, config.Get.Folder, config.Get.URLPrefix)
		} else {
			handler = handle.IgnoreIndex(handler)
		}
	}
	// If configured, apply wildcard CORS support.
	if config.Get.Cors {
		handler = handle.AddCorsWildcardHeaders(handler)
	}

	// If configured, apply key code access control.
	if "" != config.Get.AccessKey {
		handler = handle.AddAccessKey(handler, config.Get.AccessKey)
	}

	return
}

// listenerSelector returns the appropriate listener handler based on
// configuration.
func listenerSelector() (listener handle.ListenerFunc) {
	// Serve files over HTTP or HTTPS based on paths to TLS files being
	// provided.
	if 0 < len(config.Get.TLSCert) {
		handle.SetMinimumTLSVersion(config.Get.TLSMinVers)
		listener = handle.TLSListening(
			config.Get.TLSCert,
			config.Get.TLSKey,
		)
	} else {
		listener = handle.Listening()
	}
	return
}
