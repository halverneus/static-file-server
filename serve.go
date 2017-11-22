package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	version = "Version 1.1"

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
	folder := env("FOLDER", "/web") + "/"
	host := env("HOST", "")
	port := env("PORT", "8080")
	tlsCert := env("TLS_CERT", "")
	tlsKey := env("TLS_KEY", "")
	urlPrefix := env("URL_PREFIX", "")

	// If HTTPS is to be used, verify both TLS_* environment variables are set.
	if 0 < len(tlsCert) || 0 < len(tlsKey) {
		if 0 == len(tlsCert) || 0 == len(tlsKey) {
			log.Fatalln(
				"If value for environment variable 'TLS_CERT' or 'TLS_KEY' is set " +
					"then value for environment variable 'TLS_KEY' or 'TLS_CERT' must " +
					"also be set.",
			)
		}
	}

	// If the URL path prefix is to be used, verify it is properly formatted.
	if 0 < len(urlPrefix) &&
		(!strings.HasPrefix(urlPrefix, "/") || strings.HasSuffix(urlPrefix, "/")) {
		log.Fatalln(
			"Value for environment variable 'URL_PREFIX' must start " +
				"with '/' and not end with '/'. Example: '/my/prefix'",
		)
	}

	// Choose and set the appropriate, optimized static file serving function.
	var handler http.HandlerFunc
	if 0 == len(urlPrefix) {
		handler = basicHandler(folder)
	} else {
		handler = prefixHandler(folder, urlPrefix)
	}
	http.HandleFunc("/", handler)

	// Serve files over HTTP or HTTPS based on paths to TLS files being provided.
	if 0 == len(tlsCert) {
		log.Fatalln(http.ListenAndServe(host+":"+port, nil))
	} else {
		log.Fatalln(http.ListenAndServeTLS(host+":"+port, tlsCert, tlsKey, nil))
	}
}

// basicHandler serves files from the folder passed.
func basicHandler(folder string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, folder+r.URL.Path)
	}
}

// prefixHandler removes the URL path prefix before serving files from the
// folder passed.
func prefixHandler(folder, urlPrefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, urlPrefix) {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, folder+strings.TrimPrefix(r.URL.Path, urlPrefix))
	}
}

// env returns the value for an environment variable or, if not set, a fallback
// value.
func env(key, fallback string) string {
	if value := os.Getenv(key); 0 < len(value) {
		return value
	}
	return fallback
}
