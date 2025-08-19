// Package message is a simple package that provides a message interface.
package message

import (
	"encoding/json"
	"fmt"
	html "html/template"
	"strings"
	text "text/template"
	"time"

	"go.yaml.in/yaml/v2"

	"github.com/aide-family/magicbox/strutil"
)

type Message interface {
	Message() []byte
}

func TextFormatter(format string, data any) (string, error) {
	if format == "" {
		return "", fmt.Errorf("format is null")
	}
	if data == nil {
		return "", fmt.Errorf("data is nil")
	}

	t, err := text.New("text/template").Funcs(templateFuncMap).Parse(format)
	if err != nil {
		return "", err
	}
	tpl := text.Must(t, err)
	resultIoWriter := new(strings.Builder)

	if err = tpl.Execute(resultIoWriter, data); err != nil {
		return "", err
	}
	return resultIoWriter.String(), nil
}

func HtmlFormatter(format string, data any) (string, error) {
	if format == "" {
		return "", fmt.Errorf("format is null")
	}
	if data == nil {
		return "", fmt.Errorf("data is nil")
	}

	t, err := html.New("html/template").Funcs(templateFuncMap).Parse(format)
	if err != nil {
		return "", err
	}
	tpl := html.Must(t, err)
	resultIoWriter := new(strings.Builder)

	if err = tpl.Execute(resultIoWriter, data); err != nil {
		return "", err
	}
	return resultIoWriter.String(), nil
}

var templateFuncMap = map[string]any{
	"now":          time.Now,
	"hasPrefix":    strings.HasPrefix,
	"hasSuffix":    strings.HasSuffix,
	"contains":     strings.Contains,
	"trimSpace":    strings.TrimSpace,
	"trimPrefix":   strings.TrimPrefix,
	"trimSuffix":   strings.TrimSuffix,
	"toUpper":      strings.ToUpper,
	"toLower":      strings.ToLower,
	"replace":      strings.Replace,
	"split":        strings.Split,
	"mask":         strutil.MaskString,
	"maskEmail":    strutil.MaskEmail,
	"maskPhone":    strutil.MaskPhone,
	"maskBankCard": strutil.MaskBankCard,
	"title":        strutil.Title,
	"json":         json.Marshal,
	"yaml":         yaml.Marshal,
}
