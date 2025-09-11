package serialize_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aide-family/magicbox/serialize"
	"github.com/stretchr/testify/assert"
)

func TestTOMLMarshal(t *testing.T) {
	type Test struct {
		Foo string `toml:"foo"`
	}

	toml, err := serialize.TOMLMarshal(Test{Foo: "bar"})
	assert.Nil(t, err)
	assert.Equal(t, strings.TrimSpace(string(toml)), `foo = 'bar'`)

	var test Test
	err = serialize.TOMLUnmarshal(toml, &test)
	assert.Nil(t, err)
	assert.Equal(t, test.Foo, "bar")

	var test2 Test
	err = serialize.TOMLDecoder(bytes.NewReader(toml), &test2)
	assert.Nil(t, err)
	assert.Equal(t, test2.Foo, "bar")

	buf := bytes.NewBuffer([]byte{})
	err = serialize.TOMLEncoder(buf, &test)
	assert.Nil(t, err)
	assert.Equal(t, buf.String(), string(toml))
}
