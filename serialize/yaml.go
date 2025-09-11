package serialize

import (
	"io"

	yaml "gopkg.in/yaml.v2"
)

var (
	yamlMarshal   = yaml.Marshal
	yamlUnmarshal = yaml.Unmarshal
	yamlDecoder   = yaml.NewDecoder
	yamlEncoder   = yaml.NewEncoder
)

func YAMLMarshal(v any) ([]byte, error) {
	return yamlMarshal(v)
}

func YAMLUnmarshal(data []byte, v any) error {
	return yamlUnmarshal(data, v)
}

func YAMLDecoder(r io.Reader, v any) error {
	return yamlDecoder(r).Decode(v)
}

func YAMLEncoder(w io.Writer, v any) error {
	return yamlEncoder(w).Encode(v)
}
