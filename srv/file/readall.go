package filesrv

import (
	"io"
	"os"

	"github.com/reiver/go-erorr"
)

							
//@TODO: add caching for most frequently used files
							

func ReadAll(filename string) (string, error) {
	file, err := os.Open(filename)
	if nil != err {
		var nada string
		return nada, erorr.Errorf("problem opening file %q: %w", filename, err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if nil != err {
		var nada string
		return nada, erorr.Errorf("problem reading-all the content of file %q: %w", filename, err)
	}

	return string(bytes), nil
}
