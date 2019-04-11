package help

import (
	"fmt"
)

// Run print operation.
func Run() error {
	fmt.Println(Text)
	return nil
}

var (
	// Text for directly accessing help.
	Text = `
NAME
    static-file-server

SYNOPSIS
    static-file-server
    static-file-server [ -c | -config | --config ] /path/to/config.yml
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
    DEBUG
        When set to 'true' enables additional logging, including the
        configuration used and an access log for each request. IMPORTANT NOTE:
        The configuration summary is printed to stdout while logs generated
        during execution are printed to stderr. Default value is 'false'.
    FOLDER
        The path to the folder containing the contents to be served over
        HTTP(s). If not supplied, defaults to '/web' (for Docker reasons).
    HOST
        The hostname used for binding. If not supplied, contents will be served
        to a client without regard for the hostname.
    PORT
        The port used for binding. If not supplied, defaults to port '8080'.
    REFERRERS
        A comma-separated list of acceped Referrers based on the 'Referer' HTTP
        header. If incoming header value is not in the list, a 403 HTTP error is
        returned. To accept requests without a 'Referer' HTTP header in addition
        to the whitelisted values, include an empty value (either with a leading
        comma in the environment variable or with an empty list item in the YAML
        configuration file) as demonstrated in the second example. If not
        supplied the 'Referer' HTTP header is ignored.
        Examples:
          REFERRERS='http://localhost,https://some.site,http://other.site:8080'
          REFERRERS=',http://localhost,https://some.site,http://other.site:8080'
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

CONFIGURATION FILE
    Configuration can also managed used a YAML configuration file. To select the
    configuration values using the YAML file, pass in the path to the file using
    the appropriate flags (-c, --config). Environment variables take priority
    over the configuration file. The following is an example configuration using
    the default values.

    Example config.yml with defaults:
    ----------------------------------------------------------------------------
    debug: false
    folder: /web
    host: ""
    port: 8080
    referrers: []
    show-listing: true
    tls-cert: ""
    tls-key: ""
    url-prefix: ""
    ----------------------------------------------------------------------------

    Example config.yml with possible alternative values:
    ----------------------------------------------------------------------------
    debug: true
    folder: /var/www
    port: 80
    referrers:
      - http://localhost
      - https://mydomain.com
    ----------------------------------------------------------------------------

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

        export FOLDER=/var/www
        static-file-server -c config.yml
            Result: Runs with values from config.yml, but with the folder being
                    served overridden by the FOLDER environment variable.

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
