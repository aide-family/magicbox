package load

import (
	"os"
	"path/filepath"
	"regexp"

	"buf.build/go/protoyaml"
	"github.com/joho/godotenv"
	"google.golang.org/protobuf/proto"
)

func Load(cfgPath string, bootstrap proto.Message) error {
	_ = godotenv.Load()

	// 展开 ~ 路径
	cfgPath = ExpandHomeDir(cfgPath)

	if err := walk(cfgPath, bootstrap); err != nil {
		return err
	}
	return nil
}

func walk(cfgPath string, bootstrap proto.Message) error {
	return filepath.Walk(cfgPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			yamlBytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			err = protoyaml.UnmarshalOptions{
				Path: cfgPath,
			}.Unmarshal(ResolveEnv(yamlBytes), bootstrap)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// ResolveEnv resolves environment variables in content.
// It supports ${VAR} and ${VAR:default} syntax.
// If default value is empty (${VAR:}), it will use empty string as default.
func ResolveEnv(content []byte) []byte {
	// 支持 ${VAR} 和 ${VAR:default} 和 ${VAR:}（空默认值）
	regex := regexp.MustCompile(`\$\{(\w+)(?::([^}]*))?}`)

	return regex.ReplaceAllFunc(content, func(match []byte) []byte {
		matches := regex.FindSubmatch(match)
		envKey := string(matches[1])
		var defaultValue string

		// 如果有冒号，说明有默认值部分（可能是空的）
		if len(matches) > 2 && matches[2] != nil {
			defaultValue = string(matches[2])
		}

		if value, exists := os.LookupEnv(envKey); exists {
			return []byte(value)
		}
		return []byte(defaultValue)
	})
}
