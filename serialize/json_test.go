package serialize_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aide-family/magicbox/serialize"
	"github.com/stretchr/testify/assert"
)

func TestJSONMarshal(t *testing.T) {
	type Test struct {
		Foo string `json:"foo"`
	}

	json, err := serialize.JSONMarshal(Test{Foo: "bar"})
	assert.Nil(t, err)
	assert.Equal(t, string(json), `{"foo":"bar"}`)

	var test Test
	err = serialize.JSONUnmarshal(json, &test)
	assert.Nil(t, err)
	assert.Equal(t, test.Foo, "bar")

	var test2 Test
	err = serialize.JSONDecoder(bytes.NewReader(json), &test2)
	assert.Nil(t, err)
	assert.Equal(t, test2.Foo, "bar")

	buf := bytes.NewBuffer([]byte{})
	err = serialize.JSONEncoder(buf, &test)
	assert.Nil(t, err)
	assert.Equal(t, strings.TrimSpace(buf.String()), string(json))
}
