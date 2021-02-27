package mock

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

func Fixture(fname string) io.ReadCloser {

	fpath := fmt.Sprintf("testdata/%s", fname)

	f, err := ioutil.ReadFile(fpath)
	if err != nil {
		panic(fmt.Sprintf("Cannot read fixture file %s: %s", fpath, err))
	}

	return ioutil.NopCloser(bytes.NewReader(f))
}
