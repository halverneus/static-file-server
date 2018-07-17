package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/halverneus/static-file-server/config"
	"github.com/halverneus/static-file-server/handle"
)

var (
	version = "Version 1.3"

	help = `
NAME
    static-file-server

SYNOPSIS
    static-file-server
    static-file-server [ help | -help | --help ]
    static-file-server [ version | -version | --version ]

DESCRIPTION
    The Static File Server is intended to be a tiny, fast and simple solution
    for serving files over HTTP. The features included are limited to make to
    binding to a host name and port, selecting a folder to serve, choosing a
    URL path prefix and selecting TLS certificates. If you want really awesome
    reverse proxy features, I recommend Nginx.

DEPENDENCIES
    None... not even libc!

ENVIRONMENT VARIABLES
    FOLDER
        The path to the folder containing the contents to be served over
        HTTP(s). If not supplied, defaults to '/web' (for Docker reasons).
    HOST
        The hostname used for binding. If not supplied, contents will be served
        to a client without regard for the hostname.
    PORT
        The port used for binding. If not supplied, defaults to port '8080'.
    SHOW_LISTING
        Automatically serve the index file for the directory if requested. For
        example, if the client requests 'http://127.0.0.1/' the 'index.html'
        file in the root of the directory being served is returned. If the value
        is set to 'false', the same request will return a 'NOT FOUND'. Default
        value is 'true'.
    TLS_CERT
        Path to the TLS certificate file to serve files using HTTPS. If supplied
        then TLS_KEY must also be supplied. If not supplied, contents will be
        served via HTTP.
    TLS_KEY
        Path to the TLS key file to serve files using HTTPS. If supplied then
        TLS_CERT must also be supplied. If not supplied, contents will be served
        via HTTPS
    URL_PREFIX
        The prefix to use in the URL path. If supplied, then the prefix must
        start with a forward-slash and NOT end with a forward-slash. If not
        supplied then no prefix is used.

USAGE
    FILE LAYOUT
       /var/www/sub/my.file
       /var/www/index.html

    COMMAND
        export FOLDER=/var/www/sub
        static-file-server
            Retrieve with: wget http://localhost:8080/my.file
                           wget http://my.machine:8080/my.file

        export FOLDER=/var/www
        export HOST=my.machine
        export PORT=80
        static-file-server
            Retrieve with: wget http://my.machine/sub/my.file

        export FOLDER=/var/www/sub
        export HOST=my.machine
        export PORT=80
        export URL_PREFIX=/my/stuff
        static-file-server
            Retrieve with: wget http://my.machine/my/stuff/my.file

        export FOLDER=/var/www/sub
        export TLS_CERT=/etc/server/my.machine.crt
        export TLS_KEY=/etc/server/my.machine.key
        static-file-server
            Retrieve with: wget https://my.machine:8080/my.file

        export FOLDER=/var/www/sub
        export PORT=443
        export TLS_CERT=/etc/server/my.machine.crt
        export TLS_KEY=/etc/server/my.machine.key
        static-file-server
            Retrieve with: wget https://my.machine/my.file

        export FOLDER=/var/www
        export PORT=80
        export SHOW_LISTING=true  # Default behavior
        static-file-server
            Retrieve 'index.html' with: wget http://my.machine/

        export FOLDER=/var/www
        export PORT=80
        export SHOW_LISTING=false
        static-file-server
            Returns 'NOT FOUND': wget http://my.machine/
`
)

func main() {
	// Evaluate and execute subcommand if supplied.
	if 1 < len(os.Args) {
		arg := os.Args[1]
		switch {
		case strings.Contains(arg, "help"):
			fmt.Println(help)
		case strings.Contains(arg, "version"):
			fmt.Println(version)
		default:
			name := os.Args[0]
			log.Fatalf("Unknown argument: %s. Try '%s help'.", arg, name)
		}
		return
	}

	// Collect environment variables.
	if err := config.Load(""); nil != err {
		log.Fatalf("While loading configuration got %v", err)
	}

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
	log.Fatalln(listener(binding, handler))
}
