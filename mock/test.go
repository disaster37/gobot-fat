package mock

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func Fixture(fname string) io.ReadCloser {

	fpath := fmt.Sprintf("testdata/%s", fname)

	f, err := os.ReadFile(fpath)
	if err != nil {
		panic(fmt.Sprintf("Cannot read fixture file %s: %s", fpath, err))
	}

	return io.NopCloser(bytes.NewReader(f))
}
