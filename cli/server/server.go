package server

import (
	"fmt"
	"net/http"

	"github.com/halverneus/static-file-server/config"
	"github.com/halverneus/static-file-server/handle"
)

// Run server.
func Run() error {
	// Choose and set the appropriate, optimized static file serving function.
	var handler http.HandlerFunc
	if 0 == len(config.Get.URLPrefix) {
		handler = handle.Basic(config.Get.Folder)
	} else {
		handler = handle.Prefix(config.Get.Folder, config.Get.URLPrefix)
	}

	// Determine whether index files should hidden.
	if !config.Get.ShowListing {
		handler = handle.IgnoreIndex(handler)
	}

	// Serve files over HTTP or HTTPS based on paths to TLS files being provided.
	var listener handle.ListenerFunc
	if 0 < len(config.Get.TLSCert) {
		listener = handle.TLSListening(
			config.Get.TLSCert,
			config.Get.TLSKey,
		)
	} else {
		listener = handle.Listening()
	}

	binding := fmt.Sprintf("%s:%d", config.Get.Host, config.Get.Port)
	return listener(binding, handler)
}
