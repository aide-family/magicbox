package serialize

import (
	"encoding/json"
	"io"

	"github.com/aide-family/magicbox/pointer"
)

var (
	jsonMarshal   = json.Marshal
	jsonUnmarshal = json.Unmarshal
	jsonDecoder   = json.NewDecoder
	jsonEncoder   = json.NewEncoder
)

func RegisterJSONMarshal(f func(v any) ([]byte, error)) {
	jsonMarshal = f
}

func RegisterJSONUnmarshal(f func(data []byte, v any) error) {
	jsonUnmarshal = f
}

func JSONMarshal(v any) ([]byte, error) {
	if pointer.IsNil(v) {
		return nil, nil
	}
	return jsonMarshal(v)
}

func JSONUnmarshal(data []byte, v any) error {
	if pointer.IsNil(v) {
		return nil
	}
	return jsonUnmarshal(data, v)
}

func JSONDecoder(r io.Reader, v any) error {
	return jsonDecoder(r).Decode(v)
}

func JSONEncoder(w io.Writer, v any) error {
	return jsonEncoder(w).Encode(v)
}
