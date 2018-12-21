package version

import (
	"fmt"
	"runtime"
)

// Run print operation.
func Run() error {
	fmt.Printf("%s\n%s\n", VersionText, GoVersionText)
	return nil
}

var (
	// version is the application version set during build.
	version string

	// VersionText for directly accessing the static-file-server version.
	VersionText = fmt.Sprintf("v%s", version)

	// GoVersionText for directly accessing the version of the Go runtime
	// compiled with the static-file-server.
	GoVersionText = runtime.Version()
)
