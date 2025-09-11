package serialize

import (
	"io"

	"github.com/pelletier/go-toml/v2"
)

var (
	tomlMarshal   = toml.Marshal
	tomlUnmarshal = toml.Unmarshal
	tomlDecoder   = toml.NewDecoder
	tomlEncoder   = toml.NewEncoder
)

func TOMLMarshal(v any) ([]byte, error) {
	return tomlMarshal(v)
}

func TOMLUnmarshal(data []byte, v any) error {
	return tomlUnmarshal(data, v)
}

func TOMLDecoder(r io.Reader, v any) error {
	return tomlDecoder(r).Decode(v)
}

func TOMLEncoder(w io.Writer, v any) error {
	return tomlEncoder(w).Encode(v)
}
