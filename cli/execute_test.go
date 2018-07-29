package cli

import (
	"errors"
	"flag"
	"os"
	"testing"
)

func TestSetupFlags(t *testing.T) {
	app := os.Args[0]

	file := "file.txt"
	wConfig := "Config (file.txt)"

	testCases := []struct {
		name    string
		args    []string
		config  string
		help    bool
		version bool
	}{
		{"Empty args", []string{app}, "", false, false},
		{"Help (--help)", []string{app, "--help"}, "", true, false},
		{"Help (-help)", []string{app, "-help"}, "", true, false},
		{"Help (-h)", []string{app, "-h"}, "", true, false},
		{"Version (--version)", []string{app, "--version"}, "", false, true},
		{"Version (-version)", []string{app, "-version"}, "", false, true},
		{"Version (-v)", []string{app, "-v"}, "", false, true},
		{"Config ()", []string{app, "--config", ""}, "", false, false},
		{wConfig, []string{app, "--config", file}, file, false, false},
		{wConfig, []string{app, "--config=file.txt"}, file, false, false},
		{wConfig, []string{app, "-config", file}, file, false, false},
		{wConfig, []string{app, "-config=file.txt"}, file, false, false},
		{wConfig, []string{app, "-c", file}, file, false, false},
		{"All set", []string{app, "-h", "-v", "-c", file}, file, true, true},
	}

	reset := func() {
		option.configFile = ""
		option.helpFlag = false
		option.versionFlag = false
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reset()
			os.Args = tc.args
			flag.Parse()

			if option.configFile != tc.config {
				t.Errorf(
					"For options [%v] expected a config file of %s but got %s",
					tc.args, tc.config, option.configFile,
				)
			}
			if option.helpFlag != tc.help {
				t.Errorf(
					"For options [%v] expected help flag of %t but got %t",
					tc.args, tc.help, option.helpFlag,
				)
			}
			if option.versionFlag != tc.version {
				t.Errorf(
					"For options [%v] expected version flag of %t but got %t",
					tc.args, tc.version, option.versionFlag,
				)
			}
		})
	}
}

func TestExecuteAndSelection(t *testing.T) {
	app := os.Args[0]

	runHelpFuncError := errors.New("help")
	runHelpFunc = func() error {
		return runHelpFuncError
	}
	runVersionFuncError := errors.New("version")
	runVersionFunc = func() error {
		return runVersionFuncError
	}
	runServerFuncError := errors.New("server")
	runServerFunc = func() error {
		return runServerFuncError
	}
	unknownArgsFuncError := errors.New("unknown")
	unknownArgsFunc = func(Args) func() error {
		return func() error {
			return unknownArgsFuncError
		}
	}

	reset := func() {
		option.configFile = ""
		option.helpFlag = false
		option.versionFlag = false
	}

	testCases := []struct {
		name   string
		args   []string
		result error
	}{
		{"Help", []string{app, "help"}, runHelpFuncError},
		{"Help", []string{app, "--help"}, runHelpFuncError},
		{"Version", []string{app, "version"}, runVersionFuncError},
		{"Version", []string{app, "--version"}, runVersionFuncError},
		{"Serve", []string{app}, runServerFuncError},
		{"Unknown", []string{app, "unknown"}, unknownArgsFuncError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reset()
			os.Args = tc.args

			if err := Execute(); tc.result != err {
				t.Errorf(
					"Expected error for %v but got %v",
					tc.result, err,
				)
			}
		})
	}
}

func TestUnknownArgs(t *testing.T) {
	errFunc := unknownArgs(Args{"unknown"})
	if err := errFunc(); nil == err {
		t.Errorf(
			"Expected a given unknown argument error but got %v",
			err,
		)
	}
}

func TestWithConfig(t *testing.T) {
	configError := errors.New("config")
	routineError := errors.New("routine")
	routine := func() error { return routineError }

	testCases := []struct {
		name       string
		loadConfig func(string) error
		result     error
	}{
		{"Config error", func(string) error { return configError }, configError},
		{"Routine error", func(string) error { return nil }, routineError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loadConfig = tc.loadConfig
			errFunc := withConfig(routine)
			if err := errFunc(); tc.result != err {
				t.Errorf("Expected error %v but got %v", tc.result, err)
			}
		})
	}
}
