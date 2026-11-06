// Package status parses the ONT Status.json endpoint into a varid → varvalue map.
package status

import (
	"encoding/json"
	"io"
)

type statusVar struct {
	VarType  string `json:"vartype"`
	VarID    string `json:"varid"`
	VarValue string `json:"varvalue"`
}

// Parse builds a varid -> varvalue map from the ONT status JSON array.
func Parse(r io.Reader) (map[string]string, error) {
	dec := json.NewDecoder(r)

	var vars []statusVar
	if err := dec.Decode(&vars); err != nil {
		return nil, err
	}

	out := make(map[string]string, len(vars))
	for _, v := range vars {
		if v.VarID == "" {
			continue
		}

		out[v.VarID] = v.VarValue
	}

	return out, nil
}
