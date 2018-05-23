package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/scrypt"

	"github.com/howeyc/gopass"
)

var (
	version = "Version 1.2"

	authHelp = `
NAME
    static-file-server auth

SYNOPSIS
    static-file-server auth [ help | -help | --help ]
    static-file-server auth list
    static-file-server auth add $username [ $password ]
    static-file-server auth update $username [ $password ]
    static-file-server auth remove $username

DESCRIPTION
    The Static File Server authentication sub-command is used to securely modify
    a credential file for use with basic authentication.

DEPENDENCIES
    None... not even libc!

ENVIRONMENT VARIABLES
    CREDENTIALS
        The path to a file that contains valid credentials, or the path to a
        file that will have credentials modified. By having this variable set
        basic authentication will automatically be used. Credentials must be
        added to the file using the 'auth add' command prior to use. If no
        credentials are added, then basic authentication will always fail (fails
        secure). It is HIGHLY RECOMMENDED to use this with TLS certificates. If
        you are not using TLS certificates, don't use credentials that are
        important. Username and password are case-sensitive. Variable must be
        set to use any authentication sub-commands, with the exception of
        requesting help.

COMMANDS
    add $username [ $password ]
        Add a new credential to the credential file. If the username is not set
        or is the same as an existing username, the command will fail. If the
        password for a user needs to be updated, use 'auth update'. If the
        password is not supplied, it will be requested during execution.
        ARGS:
            $username: The case-sensitive username to be added.
            $password: The case-sensitive password to be associated with the new
                       username.

    help
        Prints this help documentation.

    list
        List all usernames in the credential file.

    remove $username
        Remove an existing username from the credential file. If the username is
        not set or doesn't match an existing username, the command will fail.
        ARGS:
            $username: The case-sensitive username to be removed.

    update $username [ $password ]
        Update an existing username with a new password in the credential file.
        If the username is not set or doesn't match an existing username, the
        command will fail. If a new user needs to be added, use 'auth add'. If
        the password is not supplied, it will be requested during execution.
        ARGS:
            $username: The case-sensitive username to be updated with a new
                       password.
            $password: The case-sensitive passwrod to be associated with the
                       existing username.
`

	help = `
NAME
    static-file-server

SYNOPSIS
    static-file-server
    static-file-server [ help | -help | --help ]
    static-file-server [ version | -version | --version ]
    static-file-server auth [ help | -help | --help ]

DESCRIPTION
    The Static File Server is intended to be a tiny, fast and simple solution
    for serving files over HTTP. The features included are limited to make to
    binding to a host name and port, selecting a folder to serve, choosing a
    URL path prefix and selecting TLS certificates. If you want really awesome
    reverse proxy features, I recommend Nginx.

DEPENDENCIES
    None... not even libc!

ENVIRONMENT VARIABLES
    CREDENTIALS
        The path to a file that contains valid credentials, or the path to a
        file that will have credentials modified. By having this variable set
        basic authentication will automatically be used. Credentials must be
        added to the file using the 'auth add' command prior to use. If no
        credentials are added, then basic authentication will always fail (fails
        secure). It is HIGHLY RECOMMENDED to use this with TLS certificates. If
        you are not using TLS certificates, don't use credentials that are
        important. Username and password are case-sensitive.
    FAST_AUTH
        It is recommended to use CREDENTIALS and not FAST_AUTH. FAST_AUTH is
        only provided for users that want to share files for a short time only
        using an unimportant password. If CREDENTIALS is set, FAST_AUTH will be
        ignored. The value of FAST_AUTH is a colon (:) delimited username and
        password (for example: FAST_AUTH=user:password). Username and password
        are case-sensitive.
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

COMMANDS
    [No commands supplied]
        Serve static files based on passed environment variables. Server will
        continue to run until shutdown of the service is requested.

    auth
        The CREDENTIALS environment variable must be set to use the
        authorization command. The authorization command is used to securely
        add, modify and remove user credentials for basic authentication. For
        more information, run 'static-file-server auth help'.

    help
        Prints this help documentation.
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

        export CREDENTIALS=credentials.json
        export FOLDER=/var/www/sub
        export PORT=443
        export TLS_CERT=/etc/server/my.machine.crt
        export TLS_KEY=/etc/server/my.machine.key
        static-file-server auth add 'user' 'pass'
        static-file-server auth add 'john' '12345'
        static-file-server
            Retrieve with:
                wget --user 'user' --password 'pass' --auth-no-challenge \
                    https://my.machine/my.file
                wget --user 'john' --password '12345' --auth-no-challenge \
                    https://my.machine/my.file

        export FAST_AUTH=user:pass
        export FOLDER=/var/www/sub
        export PORT=443
        export TLS_CERT=/etc/server/my.machine.crt
        export TLS_KEY=/etc/server/my.machine.key
        static-file-server
            Retrieve with:
                wget --user 'user' --password 'pass' --auth-no-challenge \
                    https://my.machine/my.file
`

	credentials map[string]*credential

	envvar struct {
		credentialFile string
		fastAuth       string
		folder         string
		host           string
		port           string
		showListing    bool
		tlsCert        string
		tlsKey         string
		urlPrefix      string
	}
)

// credential that is stored and used for authentication.
type credential struct {
	// Salt (random) for uniquely encrypting each password.
	Salt string `json:"salt"`
	// Password after it has been encrypted with the salt.
	Password string `json:"pass"`
}

// matches the passed uncrypted password against the assigned encrypted
// password. Returns true if the passwords match.
func (c *credential) matches(password string) (matches bool, err error) {
	matches = false
	var enc string
	if enc, err = c.encrypt(password, c.Salt); nil != err {
		return
	}
	matches = enc == c.Password
	return
}

// update the encrypted credential password with a new, unencrypted password.
func (c *credential) update(password string) (err error) {
	// Create a new salt value.
	rawSalt := make([]byte, 24)
	var n int
	if n, err = io.ReadFull(rand.Reader, rawSalt); nil != err {
		return
	}
	if len(rawSalt) != n {
		err = errors.New("failed to create random password salt")
		return
	}

	// Store salt and encrypted password.
	c.Salt = base64.StdEncoding.EncodeToString(rawSalt)
	c.Password, err = c.encrypt(password, c.Salt)
	return
}

// encrypt a password/salt combination to the encrypted password.
func (c *credential) encrypt(password, salt string) (enc string, err error) {
	// Decode salt from Base64 to raw bytes.
	var rawSalt []byte
	if rawSalt, err = base64.StdEncoding.DecodeString(salt); nil != err {
		return
	}

	// Encrypt password to bytes. Encryption values determined by current
	// cryptography suggestion.
	var rawEnc []byte
	if rawEnc, err = scrypt.Key(
		[]byte(password), rawSalt, 16384, 8, 1, 32,
	); nil != err {
		return
	}

	// Encode as Base64.
	enc = base64.StdEncoding.EncodeToString(rawEnc)
	return
}

func authAddMain(command string, args []string) (err error) {
	if 0 == len(args) || 2 < len(args) {
		return fmt.Errorf(
			"wrong number of arguments supplied to: '%s add', try '%s help'",
			command, command,
		)
	}
	username := args[0]
	var password string
	if 1 < len(args) {
		password = args[1]
	}

	if err = loadCredentials(envvar.credentialFile); nil != err {
		fmt.Printf(
			"WARNING: Credential file '%s' doesn't already exist... creating.\n",
			envvar.credentialFile,
		)
		err = nil
	}
	if _, found := credentials[username]; found {
		return fmt.Errorf(
			"user '%s' already exists, use 'update' in place of 'add'",
			username,
		)
	}
	if 0 == len(password) {
		if password, err = getPassword(); nil != err {
			return
		}
	}
	cred := &credential{}
	if err = cred.update(password); nil != err {
		return
	}
	credentials[username] = cred
	return saveCredentials(envvar.credentialFile)
}

func authListMain(command string, args []string) (err error) {
	if 0 != len(args) {
		return fmt.Errorf(
			"'%s list' does not accept any arguments, try '%s help'",
			command, command,
		)
	}

	if err = loadCredentials(envvar.credentialFile); nil != err {
		return
	}
	for username := range credentials {
		fmt.Println(username)
	}
	return
}

func authRemoveMain(command string, args []string) (err error) {
	if 1 != len(args) {
		return fmt.Errorf(
			"'%s remove' requires exactly one argument, try '%s help'",
			command, command,
		)
	}

	username := args[0]
	if err = loadCredentials(envvar.credentialFile); nil != err {
		return
	}
	if _, found := credentials[username]; !found {
		return fmt.Errorf("user '%s' doesn't exist", username)
	}
	delete(credentials, username)
	return saveCredentials(envvar.credentialFile)
}

func authUpdateMain(command string, args []string) (err error) {
	if 0 == len(args) || 2 < len(args) {
		return fmt.Errorf(
			"wrong number of arguments supplied to: '%s update', try '%s help'",
			command, command,
		)
	}
	username := args[0]
	var password string
	if 1 < len(args) {
		password = args[1]
	}

	if err = loadCredentials(envvar.credentialFile); nil != err {
		return
	}
	cred, found := credentials[username]
	if !found {
		return fmt.Errorf(
			"user '%s' doesn't exist, use 'add' in place of 'update'",
			username,
		)
	}
	if 0 == len(password) {
		if password, err = getPassword(); nil != err {
			return
		}
	}
	if err = cred.update(password); nil != err {
		return
	}
	return saveCredentials(envvar.credentialFile)
}

func authMain(command string, args []string) {
	// Subcommand not supplied. Redirect to help.
	if 0 == len(args) {
		log.Fatalf(
			"no arguments supplied to: '%s', try '%s help'",
			command, command,
		)
	}

	subCommand := args[0]
	args = args[1:]

	if strings.Contains(subCommand, "help") {
		fmt.Println(authHelp)
		return
	}

	if 0 == len(envvar.credentialFile) {
		log.Fatalln("credential file required but not set")
	}

	var err error
	switch subCommand {
	case "add":
		err = authAddMain(command, args)
	case "list":
		err = authListMain(command, args)
	case "remove":
		err = authRemoveMain(command, args)
	case "update":
		err = authUpdateMain(command, args)
	default:
		err = fmt.Errorf(
			"unrecognized command '%s %s', try '%s help'",
			command, subCommand, command,
		)
	}
	if nil != err {
		log.Fatalln(err)
	}
}

func main() {

	// Collect environment variables.
	envvar.credentialFile = env("CREDENTIALS", "")
	envvar.fastAuth = env("FAST_AUTH", "")
	envvar.folder = env("FOLDER", "/web") + "/"
	envvar.host = env("HOST", "")
	envvar.port = env("PORT", "8080")
	envvar.showListing = envAsBool("SHOW_LISTING", true)
	envvar.tlsCert = env("TLS_CERT", "")
	envvar.tlsKey = env("TLS_KEY", "")
	envvar.urlPrefix = env("URL_PREFIX", "")

	// Evaluate and execute subcommand if supplied.
	appName := os.Args[0]
	if 1 < len(os.Args) {
		arg := os.Args[1]
		command := fmt.Sprintf("%s %s", appName, arg)
		switch {
		case strings.Contains(arg, "help"):
			fmt.Println(help)
		case strings.Contains(arg, "version"):
			fmt.Println(version)
		case "auth" == arg:
			authMain(command, os.Args[2:])
		default:
			log.Fatalf("Unknown argument: %s. Try '%s help'.", arg, appName)
		}
		return
	}

	// If HTTPS is to be used, verify both TLS_* environment variables are set.
	if 0 < len(envvar.tlsCert) || 0 < len(envvar.tlsKey) {
		if 0 == len(envvar.tlsCert) || 0 == len(envvar.tlsKey) {
			log.Fatalln(
				"If value for environment variable 'TLS_CERT' or 'TLS_KEY' is set " +
					"then value for environment variable 'TLS_KEY' or 'TLS_CERT' must " +
					"also be set.",
			)
		}
	}

	// If the URL path prefix is to be used, verify it is properly formatted.
	if 0 < len(envvar.urlPrefix) &&
		(!strings.HasPrefix(envvar.urlPrefix, "/") || strings.HasSuffix(envvar.urlPrefix, "/")) {
		log.Fatalln(
			"Value for environment variable 'URL_PREFIX' must start " +
				"with '/' and not end with '/'. Example: '/my/prefix'",
		)
	}

	// Determine whether basic authentication is needed.
	auth := func(handler http.HandlerFunc) http.HandlerFunc {
		return handler
	}
	if 0 < len(envvar.credentialFile) || 0 < len(envvar.fastAuth) {
		if 0 < len(envvar.credentialFile) {
			if err := loadCredentials(envvar.credentialFile); nil != err {
				log.Fatalln(err)
			}
		} else {
			credentials = make(map[string]*credential)
			parts := strings.Split(envvar.fastAuth, ":")
			if 2 != len(parts) || 0 == len(parts[0]) || 0 == len(parts[1]) {
				log.Fatalln(
					"'FAST_AUTH' must have exactly one colon (:) to separate " +
						"the username from the password (username:password)",
				)
			}
			cred := &credential{}
			if err := cred.update(parts[1]); nil != err {
				log.Fatalln(err)
			}
			credentials[parts[0]] = cred
		}
		auth = func(handler http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				defer r.Body.Close()
				username, password, ok := r.BasicAuth()
				if !ok || 0 == len(username) || 0 == len(password) {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				cred, ok := credentials[username]
				if !ok {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				allowed, err := cred.matches(password)
				if nil != err {
					log.Println(err)
					w.WriteHeader(http.StatusForbidden)
					return
				}
				if !allowed {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				handler(w, r)
			}
		}
	}

	// Choose and set the appropriate, optimized static file serving function.
	var handler http.HandlerFunc
	if 0 == len(envvar.urlPrefix) {
		handler = handleListing(envvar.showListing, basicHandler(envvar.folder))
	} else {
		handler = handleListing(envvar.showListing, prefixHandler(envvar.folder, envvar.urlPrefix))
	}
	http.HandleFunc("/", auth(handler))

	// Serve files over HTTP or HTTPS based on paths to TLS files being provided.
	if 0 == len(envvar.tlsCert) {
		log.Fatalln(http.ListenAndServe(envvar.host+":"+envvar.port, nil))
	} else {
		log.Fatalln(http.ListenAndServeTLS(envvar.host+":"+envvar.port, envvar.tlsCert, envvar.tlsKey, nil))
	}
}

// handleListing wraps an HTTP request. In the event of a folder root request,
// setting 'show' to false will automatically return 'NOT FOUND' whereas true
// will attempt to retrieve the index file of that directory.
func handleListing(show bool, serve http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !show && strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		serve(w, r)
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

// loadCredentials from an existing credentials file.
func loadCredentials(filename string) (err error) {
	credentials = make(map[string]*credential)
	if 0 == len(filename) {
		return errors.New("credential file name not set but is required")
	}
	var contents []byte
	if contents, err = ioutil.ReadFile(filename); nil != err {
		return
	}
	return json.Unmarshal(contents, &credentials)
}

// saveCredentials to the specified credentials file.
func saveCredentials(filename string) (err error) {
	if 0 == len(filename) {
		return errors.New("credential file name not set but is required")
	}
	var contents []byte
	if contents, err = json.Marshal(&credentials); nil != err {
		return
	}
	return ioutil.WriteFile(filename, contents, 0600)
}

// getPassword from the user via the terminal. Mask all characters.
func getPassword() (password string, err error) {
	maskInput := true
	var rawPassword []byte
	var rawConfirmPassword []byte
	if rawPassword, err = gopass.GetPasswdPrompt(
		"New password:", maskInput, os.Stdin, os.Stdout,
	); nil != err {
		return
	}
	if 0 == len(rawPassword) {
		err = errors.New("password may not be empty")
		return
	}
	if rawConfirmPassword, err = gopass.GetPasswdPrompt(
		"Confirm password:", maskInput, os.Stdin, os.Stdout,
	); nil != err {
		return
	}
	if !bytes.Equal(rawPassword, rawConfirmPassword) {
		err = errors.New("passwords do not match")
		return
	}
	password = string(rawPassword)
	return
}

// env returns the value for an environment variable or, if not set, a fallback
// value.
func env(key, fallback string) string {
	if value := os.Getenv(key); 0 < len(value) {
		return value
	}
	return fallback
}

// strAsBool converts the intent of the passed value into a boolean
// representation.
func strAsBool(value string) (result bool, err error) {
	lvalue := strings.ToLower(value)
	switch lvalue {
	case "0", "false", "f", "no", "n":
		result = false
	case "1", "true", "t", "yes", "y":
		result = true
	default:
		result = false
		msg := "Unknown conversion from string to bool for value '%s'"
		err = fmt.Errorf(msg, value)
	}
	return
}

// envAsBool returns the value for an environment variable or, if not set, a
// fallback value as a boolean.
func envAsBool(key string, fallback bool) bool {
	value := env(key, fmt.Sprintf("%t", fallback))
	result, err := strAsBool(value)
	if nil != err {
		log.Printf(
			"Invalid value for '%s': %v\nUsing fallback: %t",
			key, err, fallback,
		)
		return fallback
	}
	return result
}
