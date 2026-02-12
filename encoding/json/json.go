// Package json provides a JSON codec.
package json

import (
	"encoding/json"

	"github.com/aide-family/magicbox/encoding"
	kratosencoding "github.com/go-kratos/kratos/v2/encoding"
	kratosjson "github.com/go-kratos/kratos/v2/encoding/json"
)

const Name = kratosjson.Name

func init() {
	encoding.RegisterCodec(Name, &jsonCodec{
		Codec: kratosencoding.GetCodec(Name),
	})
}

type jsonCodec struct {
	kratosencoding.Codec
}

// Valid implements [encoding.Codec].
func (j *jsonCodec) Valid(data []byte) bool {
	return json.Valid(data)
}
