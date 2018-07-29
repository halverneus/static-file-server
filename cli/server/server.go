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
	// Choose and set the appropriate, optimized static file serving function.
	if 0 == len(config.Get.URLPrefix) {
		handler = handle.Basic(config.Get.Folder)
	} else {
		handler = handle.Prefix(config.Get.Folder, config.Get.URLPrefix)
	}

	// Determine whether index files should hidden.
	if !config.Get.ShowListing {
		handler = handle.IgnoreIndex(handler)
	}
	return
}

// listenerSelector returns the appropriate listener handler based on
// configuration.
func listenerSelector() (listener handle.ListenerFunc) {
	// Serve files over HTTP or HTTPS based on paths to TLS files being
	// provided.
	if 0 < len(config.Get.TLSCert) {
		listener = handle.TLSListening(
			config.Get.TLSCert,
			config.Get.TLSKey,
		)
	} else {
		listener = handle.Listening()
	}
	return
}
