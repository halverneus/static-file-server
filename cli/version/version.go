package version

import (
	"fmt"
	"runtime"
)

// Run print operation.
func Run() error {
	fmt.Printf("%s\n%s\n", Text, GoVersionText)
	return nil
}

var (
	// MajorVersion of static-file-server.
	MajorVersion = 1

	// MinorVersion of static-file-server.
	MinorVersion = 4

	// FixVersion of static-file-server.
	FixVersion = 0

	// Text for directly accessing the static-file-server version.
	Text = fmt.Sprintf(
		"v%d.%d.%d",
		MajorVersion,
		MinorVersion,
		FixVersion,
	)

	// GoVersionText for directly accessing the version of the Go runtime
	// compiled with the static-file-server.
	GoVersionText = runtime.Version()
)
