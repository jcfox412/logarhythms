package testing

import (
	"os"
	"path"
	"runtime"
)

// This function allows all tests to be run from the root of the project,
// allowing all tests access easy access to all asset files.
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
