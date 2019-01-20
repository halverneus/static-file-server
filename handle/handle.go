package handle

import (
	"log"
	"net/http"
	"strings"
)

var (
	// These assignments are for unit testing.
	listenAndServe    = http.ListenAndServe
	listenAndServeTLS = http.ListenAndServeTLS
	setHandler        = http.HandleFunc
)

var (
	server http.Server
)

// ListenerFunc accepts the {hostname:port} binding string required by HTTP
// listeners and the handler (router) function and returns any errors that
// occur.
type ListenerFunc func(string, http.HandlerFunc) error

// FileServerFunc is used to serve the file from the local file system to the
// requesting client.
type FileServerFunc func(http.ResponseWriter, *http.Request, string)

func validReferrer(s []string, e string) bool {
	if (s == nil) {
		// log.Printf("No referrers specified, all fine.")
		return true
	}

	// log.Printf("Checking referrers " + strings.Join(s, ",")  + " against " + e)

	for _, a := range s {
		// Handle blank HTTP Referer header, if configured
		if (a == "") {
			if (e == "") {
				// log.Printf("No referrer in request. Allowing.");
				return true;
			}
			// Continue loop (all strings start with "")
			continue;
		}

		// Compare header with allowed prefixes
		if strings.HasPrefix(e, a) {
			// log.Printf(strings.Join([]string{ "Referrer match", e, a }, " "));
			return true
		}
	}
	return false
}

func WithReferrers(serveFile FileServerFunc, referrers []string) FileServerFunc {
	return func(w http.ResponseWriter, r *http.Request, name string) {
		if (validReferrer(referrers, r.Referer())) {
			// log.Printf("Serving file.")
			serveFile(w, r, name)
		} else {
				// log.Printf(strings.Join([]string{"Invalid referrer", r.Referer(), "Not in", strings.Join(referrers, ",")}, " "))
				http.Error(w, strings.Join([]string{ "Invalid source", r.Referer() }, " "), 403)
				return
		}
	}
}

// WithLogging returns a function that logs information about the request prior
// to serving the requested file.
func WithLogging(serveFile FileServerFunc) FileServerFunc {
	return func(w http.ResponseWriter, r *http.Request, name string) {
		log.Printf(
			"REQ from %s: %s %s %s%s -> %s\n",
			r.Referer(),
			r.Method,
			r.Proto,
			r.Host,
			r.URL.Path,
			name,
		)
		serveFile(w, r, name)
	}
}

// Basic file handler servers files from the passed folder.
func Basic(serveFile FileServerFunc, folder string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serveFile(w, r, folder+r.URL.Path)
	}
}

// Prefix file handler is an alternative to Basic where a URL prefix is removed
// prior to serving a file (http://my.machine/prefix/file.txt will serve
// file.txt from the root of the folder being served (ignoring 'prefix')).
func Prefix(serveFile FileServerFunc, folder, urlPrefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, urlPrefix) {
			http.NotFound(w, r)
			return
		}
		serveFile(w, r, folder+strings.TrimPrefix(r.URL.Path, urlPrefix))
	}
}

// IgnoreIndex wraps an HTTP request. In the event of a folder root request,
// this function will automatically return 'NOT FOUND' as opposed to default
// behavior where the index file for that directory is retrieved.
func IgnoreIndex(serve http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		serve(w, r)
	}
}

// Listening function for serving the handler function.
func Listening() ListenerFunc {
	return func(binding string, handler http.HandlerFunc) error {
		setHandler("/", handler)
		return listenAndServe(binding, nil)
	}
}

// TLSListening function for serving the handler function with encryption.
func TLSListening(tlsCert, tlsKey string) ListenerFunc {
	return func(binding string, handler http.HandlerFunc) error {
		setHandler("/", handler)
		return listenAndServeTLS(binding, tlsCert, tlsKey, nil)
	}
}
