package util

import (
	"encoding/json"
	"io"
)

func DecodeJSONBody(body io.Reader, v interface{}) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		return err
	}
	return nil
}

// EncodeJSONBody encodes a struct into JSON and writes it to the provided io.Writer.
func EncodeJSONBody(w io.Writer, v interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ") // Optional: pretty print with indentation
	return encoder.Encode(v)
}
