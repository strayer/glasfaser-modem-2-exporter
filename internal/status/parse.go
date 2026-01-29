package status

import (
	"errors"
	"io"
)

// Parse builds a varid -> varvalue map from the ONT status JSON array.
// TODO: implement.
func Parse(_ io.Reader) (map[string]string, error) {
	return nil, errors.New("not implemented")
}
