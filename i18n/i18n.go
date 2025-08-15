// Package i18n provides a i18n service.
package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	yaml "sigs.k8s.io/yaml/goyaml.v2"
)

// New creates a new i18n bundle.
func New(config Config) (*i18n.Bundle, error) {
	lang := config.GetLang()
	newBundle := i18n.NewBundle(lang)
	format := config.GetFormat()

	switch format {
	case FormatJSON:
		newBundle.RegisterUnmarshalFunc(format.String(), json.Unmarshal)
	case FormatYAML:
		newBundle.RegisterUnmarshalFunc(format.String(), yaml.Unmarshal)
	case FormatTOML:
		newBundle.RegisterUnmarshalFunc(format.String(), toml.Unmarshal)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
	dir := config.GetDir()
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	ext := format.Ext()
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ext) {
			continue
		}
		newBundle.MustLoadMessageFile(path.Join(dir, file.Name()))
	}
	return newBundle, nil
}

// Message returns a localized message.
func Message(bundle *i18n.Bundle, lang string, key string, args ...interface{}) (string, error) {
	localize, err := i18n.NewLocalizer(bundle, lang).
		Localize(&i18n.LocalizeConfig{MessageID: key, TemplateData: args})
	if err != nil {
		return "", err
	}
	return localize, nil
}

// MessageX returns a localized message.
func MessageX(bundle *i18n.Bundle, lang string, key string, args ...interface{}) string {
	localize, _ := Message(bundle, lang, key, args...)
	return localize
}
