package serialize_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aide-family/magicbox/serialize"
	"github.com/stretchr/testify/assert"
)

func TestYAMLMarshal(t *testing.T) {
	type Test struct {
		Foo string `yaml:"foo"`
	}

	yaml, err := serialize.YAMLMarshal(Test{Foo: "bar"})
	assert.Nil(t, err)
	assert.Equal(t, strings.TrimSpace(string(yaml)), `foo: bar`)

	var test Test
	err = serialize.YAMLUnmarshal(yaml, &test)
	assert.Nil(t, err)
	assert.Equal(t, test.Foo, "bar")

	var test2 Test
	err = serialize.YAMLDecoder(bytes.NewReader(yaml), &test2)
	assert.Nil(t, err)
	assert.Equal(t, test2.Foo, "bar")

	buf := bytes.NewBuffer([]byte{})
	err = serialize.YAMLEncoder(buf, &test)
	assert.Nil(t, err)
	assert.Equal(t, buf.String(), string(yaml))
}
