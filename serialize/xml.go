package serialize

import (
	"encoding/xml"
	"io"
)

var (
	xmlMarshal   = xml.Marshal
	xmlUnmarshal = xml.Unmarshal
	xmlDecoder   = xml.NewDecoder
	xmlEncoder   = xml.NewEncoder
)

func XMLMarshal(v any) ([]byte, error) {
	return xmlMarshal(v)
}

func XMLUnmarshal(data []byte, v any) error {
	return xmlUnmarshal(data, v)
}

func XMLDecoder(r io.Reader, v any) error {
	return xmlDecoder(r).Decode(v)
}

func XMLEncoder(w io.Writer, v any) error {
	return xmlEncoder(w).Encode(v)
}
