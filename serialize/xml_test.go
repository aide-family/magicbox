package serialize_test

import (
	"bytes"
	"testing"

	"github.com/aide-family/magicbox/serialize"
	"github.com/stretchr/testify/assert"
)

func TestXMLMarshal(t *testing.T) {
	type Test struct {
		Foo string `xml:"foo"`
	}

	xml, err := serialize.XMLMarshal(Test{Foo: "bar"})
	assert.Nil(t, err)
	assert.Equal(t, string(xml), `<Test><foo>bar</foo></Test>`)

	var test Test
	err = serialize.XMLUnmarshal(xml, &test)
	assert.Nil(t, err)
	assert.Equal(t, test.Foo, "bar")

	var test2 Test
	err = serialize.XMLDecoder(bytes.NewReader(xml), &test2)
	assert.Nil(t, err)
	assert.Equal(t, test2.Foo, "bar")

	buf := bytes.NewBuffer([]byte{})
	err = serialize.XMLEncoder(buf, &test)
	assert.Nil(t, err)
	assert.Equal(t, buf.String(), string(xml))
}
