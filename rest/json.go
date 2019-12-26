package rest

import (
	"encoding/json"
	"io"
)

func readJson(r io.Reader, data interface{}) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(&data)
	return err
}

func writeJson(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(data)
}
