// Package yaml provides a YAML codec.
package yaml

import (
	"gopkg.in/yaml.v3"

	"github.com/aide-family/magicbox/encoding"
	kratosencoding "github.com/go-kratos/kratos/v2/encoding"
	kratosyaml "github.com/go-kratos/kratos/v2/encoding/yaml"
)

const Name = kratosyaml.Name

func init() {
	encoding.RegisterCodec(Name, &yamlCodec{
		Codec: kratosencoding.GetCodec(Name),
	})
}

type yamlCodec struct {
	kratosencoding.Codec
}

// Valid implements [encoding.Codec].
func (y *yamlCodec) Valid(data []byte) bool {
	return yaml.Unmarshal(data, &yaml.Node{}) == nil
}
