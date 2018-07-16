package handle

import (
	"net/http"
	"strings"
)

type ListenerFunc func(string, http.HandlerFunc) error

// Basic file handler servers files from the passed folder.
func Basic(folder string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, folder+r.URL.Path)
	}
}

// Prefix file handler is an alternative to Basic where a URL prefix is removed
// prior to serving a file (http://my.machine/prefix/file.txt will serve
// file.txt from the root of the folder being served (ignoring 'prefix')).
func Prefix(folder, urlPrefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, urlPrefix) {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, folder+strings.TrimPrefix(r.URL.Path, urlPrefix))
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
		http.HandleFunc("/", handler)
		return http.ListenAndServe(binding, nil)
	}
}

// TLSListening function for serving the handler function with encryption.
func TLSListening(tlsCert, tlsKey string) ListenerFunc {
	return func(binding string, handler http.HandlerFunc) error {
		http.HandleFunc("/", handler)
		return http.ListenAndServeTLS(binding, tlsCert, tlsKey, nil)
	}
}
