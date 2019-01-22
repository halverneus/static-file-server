package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var (
	// Get the desired configuration value.
	Get struct {
		Debug       bool     `yaml:"debug"`
		Folder      string   `yaml:"folder"`
		Host        string   `yaml:"host"`
		Port        uint16   `yaml:"port"`
		ShowListing bool     `yaml:"show-listing"`
		TLSCert     string   `yaml:"tls-cert"`
		TLSKey      string   `yaml:"tls-key"`
		URLPrefix   string   `yaml:"url-prefix"`
		Referrers   []string `yaml:"referrers"`
	}
)

const (
	debugKey       = "DEBUG"
	folderKey      = "FOLDER"
	hostKey        = "HOST"
	portKey        = "PORT"
	referrersKey   = "REFERRERS"
	showListingKey = "SHOW_LISTING"
	tlsCertKey     = "TLS_CERT"
	tlsKeyKey      = "TLS_KEY"
	urlPrefixKey   = "URL_PREFIX"
)

var (
	defaultDebug       = false
	defaultFolder      = "/web"
	defaultHost        = ""
	defaultPort        = uint16(8080)
	defaultReferrers   = []string{}
	defaultShowListing = true
	defaultTLSCert     = ""
	defaultTLSKey      = ""
	defaultURLPrefix   = ""
)

func init() {
	// init calls setDefaults to better support testing.
	setDefaults()
}

func setDefaults() {
	Get.Debug = defaultDebug
	Get.Folder = defaultFolder
	Get.Host = defaultHost
	Get.Port = defaultPort
	Get.Referrers = defaultReferrers
	Get.ShowListing = defaultShowListing
	Get.TLSCert = defaultTLSCert
	Get.TLSKey = defaultTLSKey
	Get.URLPrefix = defaultURLPrefix
}

// Load the configuration file.
func Load(filename string) (err error) {
	// If no filename provided, assign envvars.
	if "" == filename {
		overrideWithEnvVars()
		return
	}

	// Read contents from configuration file.
	var contents []byte
	if contents, err = ioutil.ReadFile(filename); nil != err {
		return
	}

	// Parse contents into 'Get' configuration.
	if err = yaml.Unmarshal(contents, &Get); nil != err {
		return
	}

	overrideWithEnvVars()
	return validate()
}

// Log the current configuration.
func Log() {
	// YAML marshalling should never error, but if it could, the result is that
	// the contents of the configuration are not logged.
	contents, _ := yaml.Marshal(&Get)

	// Log the configuration.
	fmt.Println("Using the following configuration:")
	fmt.Println(string(contents))
}

// overrideWithEnvVars the default values and the configuration file values.
func overrideWithEnvVars() {
	// Assign envvars, if set.
	Get.Debug = envAsBool(debugKey, Get.Debug)
	Get.Folder = envAsStr(folderKey, Get.Folder)
	Get.Host = envAsStr(hostKey, Get.Host)
	Get.Port = envAsUint16(portKey, Get.Port)
	Get.ShowListing = envAsBool(showListingKey, Get.ShowListing)
	Get.TLSCert = envAsStr(tlsCertKey, Get.TLSCert)
	Get.TLSKey = envAsStr(tlsKeyKey, Get.TLSKey)
	Get.URLPrefix = envAsStr(urlPrefixKey, Get.URLPrefix)
	Get.Referrers = envAsStrSlice(referrersKey, Get.Referrers)
}

// validate the configuration.
func validate() error {
	// If HTTPS is to be used, verify both TLS_* environment variables are set.
	if 0 < len(Get.TLSCert) || 0 < len(Get.TLSKey) {
		if 0 == len(Get.TLSCert) || 0 == len(Get.TLSKey) {
			msg := "if value for either 'TLS_CERT' or 'TLS_KEY' is set then " +
				"then value for the other must also be set (values are " +
				"currently '%s' and '%s', respectively)"
			return fmt.Errorf(msg, Get.TLSCert, Get.TLSKey)
		}
		if _, err := os.Stat(Get.TLSCert); nil != err {
			msg := "value of TLS_CERT is set with filename '%s' that returns %v"
			return fmt.Errorf(msg, err)
		}
		if _, err := os.Stat(Get.TLSKey); nil != err {
			msg := "value of TLS_KEY is set with filename '%s' that returns %v"
			return fmt.Errorf(msg, err)
		}
	}

	// If the URL path prefix is to be used, verify it is properly formatted.
	if 0 < len(Get.URLPrefix) &&
		(!strings.HasPrefix(Get.URLPrefix, "/") || strings.HasSuffix(Get.URLPrefix, "/")) {
		msg := "if value for 'URL_PREFIX' is set then the value must start " +
			"with '/' and not end with '/' (current value of '%s' vs valid " +
			"example of '/my/prefix'"
		return fmt.Errorf(msg, Get.URLPrefix)
	}

	return nil
}

// envAsStr returns the value of the environment variable as a string if set.
func envAsStr(key, fallback string) string {
	if value := os.Getenv(key); "" != value {
		return value
	}
	return fallback
}

// envAsStrSlice returns the value of the environment variable as a slice of
// strings if set.
func envAsStrSlice(key string, fallback []string) []string {
	if value := os.Getenv(key); "" != value {
		return strings.Split(value, ",")
	}
	return fallback
}

// envAsUint16 returns the value of the environment variable as a uint16 if set.
func envAsUint16(key string, fallback uint16) uint16 {
	// Retrieve the string value of the environment variable. If not set,
	// fallback is used.
	valueStr := os.Getenv(key)
	if "" == valueStr {
		return fallback
	}

	// Parse the string into a uint16.
	base := 10
	bitSize := 16
	valueAsUint64, err := strconv.ParseUint(valueStr, base, bitSize)
	if nil != err {
		log.Printf(
			"Invalid value for '%s': %v\nUsing fallback: %d",
			key, err, fallback,
		)
		return fallback
	}
	return uint16(valueAsUint64)
}

// envAsBool returns the value for an environment variable or, if not set, a
// fallback value as a boolean.
func envAsBool(key string, fallback bool) bool {
	// Retrieve the string value of the environment variable. If not set,
	// fallback is used.
	valueStr := os.Getenv(key)
	if "" == valueStr {
		return fallback
	}

	// Parse the string into a boolean.
	value, err := strAsBool(valueStr)
	if nil != err {
		log.Printf(
			"Invalid value for '%s': %v\nUsing fallback: %t",
			key, err, fallback,
		)
		return fallback
	}
	return value
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
