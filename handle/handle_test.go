package handle

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

var (
	baseDir             = "tmp/"
	subDir              = "sub/"
	subDeepDir          = "sub/deep/"
	tmpIndexName        = "index.html"
	tmpFileName         = "file.txt"
	tmpBadName          = "bad.txt"
	tmpSubIndexName     = "sub/index.html"
	tmpSubFileName      = "sub/file.txt"
	tmpSubBadName       = "sub/bad.txt"
	tmpSubDeepIndexName = "sub/deep/index.html"
	tmpSubDeepFileName  = "sub/deep/file.txt"
	tmpSubDeepBadName   = "sub/deep/bad.txt"

	tmpIndex        = "Space: the final frontier"
	tmpFile         = "These are the voyages of the starship Enterprise."
	tmpSubIndex     = "Its continuing mission:"
	tmpSubFile      = "To explore strange new worlds"
	tmpSubDeepIndex = "To seek out new life and new civilizations"
	tmpSubDeepFile  = "To boldly go where no one has gone before"

	nothing  = ""
	ok       = http.StatusOK
	missing  = http.StatusNotFound
	redirect = http.StatusMovedPermanently
	notFound = "404 page not found\n"

	files = map[string]string{
		baseDir + tmpIndexName:        tmpIndex,
		baseDir + tmpFileName:         tmpFile,
		baseDir + tmpSubIndexName:     tmpSubIndex,
		baseDir + tmpSubFileName:      tmpSubFile,
		baseDir + tmpSubDeepIndexName: tmpSubDeepIndex,
		baseDir + tmpSubDeepFileName:  tmpSubDeepFile,
	}

	serveFileFuncs = []FileServerFunc{
		http.ServeFile,
		WithLogging(http.ServeFile),
	}
)

func TestMain(m *testing.M) {
	code := func(m *testing.M) int {
		if err := setup(); nil != err {
			log.Fatalf("While setting up test got: %v\n", err)
		}
		defer teardown()
		return m.Run()
	}(m)
	os.Exit(code)
}

func setup() (err error) {
	for filename, contents := range files {
		if err = os.MkdirAll(path.Dir(filename), 0700); nil != err {
			return
		}
		if err = ioutil.WriteFile(
			filename,
			[]byte(contents),
			0600,
		); nil != err {
			return
		}
	}
	return
}

func teardown() (err error) {
	return os.RemoveAll("tmp")
}

func TestWithReferrers(t *testing.T) {
	forbidden := http.StatusForbidden

	ok1 := "http://valid.com"
	ok2 := "https://valid.com"
	ok3 := "http://localhost"
	bad := "http://other.pl"

	var noRefer []string
	emptyRefer := []string{}
	onlyNoRefer := []string{""}
	refer := []string{ok1, ok2, ok3}
	noWithRefer := []string{"", ok1, ok2, ok3}

	testCases := []struct {
		name   string
		refers []string
		refer  string
		code   int
	}{
		{"Nil refer list", noRefer, bad, ok},
		{"Empty refer list", emptyRefer, bad, ok},
		{"Unassigned allowed & unassigned", onlyNoRefer, "", ok},
		{"Unassigned allowed & assigned", onlyNoRefer, bad, forbidden},
		{"Whitelist with unassigned", refer, "", forbidden},
		{"Whitelist with bad", refer, bad, forbidden},
		{"Whitelist with ok1", refer, ok1, ok},
		{"Whitelist with ok2", refer, ok2, ok},
		{"Whitelist with ok3", refer, ok3, ok},
		{"Whitelist and none with unassigned", noWithRefer, "", ok},
		{"Whitelist with bad", noWithRefer, bad, forbidden},
		{"Whitelist with ok1", noWithRefer, ok1, ok},
		{"Whitelist with ok2", noWithRefer, ok2, ok},
		{"Whitelist with ok3", noWithRefer, ok3, ok},
	}

	success := func(w http.ResponseWriter, r *http.Request, name string) {
		defer r.Body.Close()
		w.WriteHeader(ok)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := WithReferrers(success, tc.refers)

			fullpath := "http://localhost/" + tmpIndexName
			req := httptest.NewRequest("GET", fullpath, nil)
			req.Header.Add("Referer", tc.refer)
			w := httptest.NewRecorder()

			handler(w, req, "")

			resp := w.Result()
			_, err := ioutil.ReadAll(resp.Body)
			if nil != err {
				t.Errorf("While reading body got %v", err)
			}
			if tc.code != resp.StatusCode {
				t.Errorf(
					"With referer '%s' in '%v' expected status code %d but got %d",
					tc.refer, tc.refers, tc.code, resp.StatusCode,
				)
			}
		})
	}
}

func TestBasicWithAndWithoutLogging(t *testing.T) {
	referer := "http://localhost"
	noReferer := ""
	testCases := []struct {
		name     string
		path     string
		code     int
		refer    string
		contents string
	}{
		{"Good base dir", "", ok, referer, tmpIndex},
		{"Good base index", tmpIndexName, redirect, referer, nothing},
		{"Good base file", tmpFileName, ok, referer, tmpFile},
		{"Bad base file", tmpBadName, missing, referer, notFound},
		{"Good subdir dir", subDir, ok, referer, tmpSubIndex},
		{"Good subdir index", tmpSubIndexName, redirect, referer, nothing},
		{"Good subdir file", tmpSubFileName, ok, referer, tmpSubFile},
		{"Good base dir", "", ok, noReferer, tmpIndex},
		{"Good base index", tmpIndexName, redirect, noReferer, nothing},
		{"Good base file", tmpFileName, ok, noReferer, tmpFile},
		{"Bad base file", tmpBadName, missing, noReferer, notFound},
		{"Good subdir dir", subDir, ok, noReferer, tmpSubIndex},
		{"Good subdir index", tmpSubIndexName, redirect, noReferer, nothing},
		{"Good subdir file", tmpSubFileName, ok, noReferer, tmpSubFile},
	}

	for _, serveFile := range serveFileFuncs {
		handler := Basic(serveFile, baseDir)
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				fullpath := "http://localhost/" + tc.path
				req := httptest.NewRequest("GET", fullpath, nil)
				req.Header.Add("Referer", tc.refer)
				w := httptest.NewRecorder()

				handler(w, req)

				resp := w.Result()
				body, err := ioutil.ReadAll(resp.Body)
				if nil != err {
					t.Errorf("While reading body got %v", err)
				}
				contents := string(body)
				if tc.code != resp.StatusCode {
					t.Errorf(
						"While retrieving %s expected status code of %d but got %d",
						fullpath, tc.code, resp.StatusCode,
					)
				}
				if tc.contents != contents {
					t.Errorf(
						"While retrieving %s expected contents '%s' but got '%s'",
						fullpath, tc.contents, contents,
					)
				}
			})
		}
	}
}

func TestPrefix(t *testing.T) {
	prefix := "/my/prefix/path/"

	testCases := []struct {
		name     string
		path     string
		code     int
		contents string
	}{
		{"Good base dir", prefix, ok, tmpIndex},
		{"Good base index", prefix + tmpIndexName, redirect, nothing},
		{"Good base file", prefix + tmpFileName, ok, tmpFile},
		{"Bad base file", prefix + tmpBadName, missing, notFound},
		{"Good subdir dir", prefix + subDir, ok, tmpSubIndex},
		{"Good subdir index", prefix + tmpSubIndexName, redirect, nothing},
		{"Good subdir file", prefix + tmpSubFileName, ok, tmpSubFile},
		{"Unknown prefix", tmpFileName, missing, notFound},
	}

	for _, serveFile := range serveFileFuncs {
		handler := Prefix(serveFile, baseDir, prefix)
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				fullpath := "http://localhost" + tc.path
				req := httptest.NewRequest("GET", fullpath, nil)
				w := httptest.NewRecorder()

				handler(w, req)

				resp := w.Result()
				body, err := ioutil.ReadAll(resp.Body)
				if nil != err {
					t.Errorf("While reading body got %v", err)
				}
				contents := string(body)
				if tc.code != resp.StatusCode {
					t.Errorf(
						"While retrieving %s expected status code of %d but got %d",
						fullpath, tc.code, resp.StatusCode,
					)
				}
				if tc.contents != contents {
					t.Errorf(
						"While retrieving %s expected contents '%s' but got '%s'",
						fullpath, tc.contents, contents,
					)
				}
			})
		}
	}
}

func TestIgnoreIndex(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		code     int
		contents string
	}{
		{"Good base dir", "", missing, notFound},
		{"Good base index", tmpIndexName, redirect, nothing},
		{"Good base file", tmpFileName, ok, tmpFile},
		{"Bad base file", tmpBadName, missing, notFound},
		{"Good subdir dir", subDir, missing, notFound},
		{"Good subdir index", tmpSubIndexName, redirect, nothing},
		{"Good subdir file", tmpSubFileName, ok, tmpSubFile},
	}

	for _, serveFile := range serveFileFuncs {
		handler := IgnoreIndex(Basic(serveFile, baseDir))
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				fullpath := "http://localhost/" + tc.path
				req := httptest.NewRequest("GET", fullpath, nil)
				w := httptest.NewRecorder()

				handler(w, req)

				resp := w.Result()
				body, err := ioutil.ReadAll(resp.Body)
				if nil != err {
					t.Errorf("While reading body got %v", err)
				}
				contents := string(body)
				if tc.code != resp.StatusCode {
					t.Errorf(
						"While retrieving %s expected status code of %d but got %d",
						fullpath, tc.code, resp.StatusCode,
					)
				}
				if tc.contents != contents {
					t.Errorf(
						"While retrieving %s expected contents '%s' but got '%s'",
						fullpath, tc.contents, contents,
					)
				}
			})
		}
	}
}

func TestListening(t *testing.T) {
	// Choose values for testing.
	called := false
	testBinding := "host:port"
	testError := errors.New("random problem")

	// Create an empty placeholder router function.
	handler := func(http.ResponseWriter, *http.Request) {}

	// Override setHandler so that multiple calls to 'http.HandleFunc' doesn't
	// panic.
	setHandler = func(string, func(http.ResponseWriter, *http.Request)) {}

	// Override listenAndServe with a function with more introspection and
	// control than 'http.ListenAndServe'.
	listenAndServe = func(
		binding string, handler http.Handler,
	) error {
		if testBinding != binding {
			t.Errorf(
				"While serving expected binding of %s but got %s",
				testBinding, binding,
			)
		}
		called = !called
		if called {
			return nil
		}
		return testError
	}

	// Perform test.
	listener := Listening()
	if err := listener(testBinding, handler); nil != err {
		t.Errorf("While serving first expected nil error but got %v", err)
	}
	if err := listener(testBinding, handler); nil == err {
		t.Errorf(
			"While serving second got nil while expecting %v", testError,
		)
	}
}

func TestTLSListening(t *testing.T) {
	// Choose values for testing.
	called := false
	testBinding := "host:port"
	testTLSCert := "test/file.pem"
	testTLSKey := "test/file.key"
	testError := errors.New("random problem")

	// Create an empty placeholder router function.
	handler := func(http.ResponseWriter, *http.Request) {}

	// Override setHandler so that multiple calls to 'http.HandleFunc' doesn't
	// panic.
	setHandler = func(string, func(http.ResponseWriter, *http.Request)) {}

	// Override listenAndServeTLS with a function with more introspection and
	// control than 'http.ListenAndServeTLS'.
	listenAndServeTLS = func(
		binding, tlsCert, tlsKey string, handler http.Handler,
	) error {
		if testBinding != binding {
			t.Errorf(
				"While serving TLS expected binding of %s but got %s",
				testBinding, binding,
			)
		}
		if testTLSCert != tlsCert {
			t.Errorf(
				"While serving TLS expected TLS cert of %s but got %s",
				testTLSCert, tlsCert,
			)
		}
		if testTLSKey != tlsKey {
			t.Errorf(
				"While serving TLS expected TLS key of %s but got %s",
				testTLSKey, tlsKey,
			)
		}
		called = !called
		if called {
			return nil
		}
		return testError
	}

	// Perform test.
	listener := TLSListening(testTLSCert, testTLSKey)
	if err := listener(testBinding, handler); nil != err {
		t.Errorf("While serving first TLS expected nil error but got %v", err)
	}
	if err := listener(testBinding, handler); nil == err {
		t.Errorf(
			"While serving second TLS got nil while expecting %v", testError,
		)
	}
}

func TestValidReferrer(t *testing.T) {
	ok1 := "http://valid.com"
	ok2 := "https://valid.com"
	ok3 := "http://localhost"
	bad := "http://other.pl"

	var noRefer []string
	emptyRefer := []string{}
	onlyNoRefer := []string{""}
	refer := []string{ok1, ok2, ok3}
	noWithRefer := []string{"", ok1, ok2, ok3}

	testCases := []struct {
		name   string
		refers []string
		refer  string
		result bool
	}{
		{"Nil refer list", noRefer, bad, true},
		{"Empty refer list", emptyRefer, bad, true},
		{"Unassigned allowed & unassigned", onlyNoRefer, "", true},
		{"Unassigned allowed & assigned", onlyNoRefer, bad, false},
		{"Whitelist with unassigned", refer, "", false},
		{"Whitelist with bad", refer, bad, false},
		{"Whitelist with ok1", refer, ok1, true},
		{"Whitelist with ok2", refer, ok2, true},
		{"Whitelist with ok3", refer, ok3, true},
		{"Whitelist and none with unassigned", noWithRefer, "", true},
		{"Whitelist with bad", noWithRefer, bad, false},
		{"Whitelist with ok1", noWithRefer, ok1, true},
		{"Whitelist with ok2", noWithRefer, ok2, true},
		{"Whitelist with ok3", noWithRefer, ok3, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validReferrer(tc.refers, tc.refer)
			if result != tc.result {
				t.Errorf(
					"With referrers of '%v' and a value of '%s' expected %t but got %t",
					tc.refers, tc.refer, tc.result, result,
				)
			}
		})
	}
}
