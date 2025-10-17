package load

import (
	"os"
	"path/filepath"
	"strings"
)

// expandHomeDir expand the ~ symbol in the path to the user's home directory
func ExpandHomeDir(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			// 如果无法获取用户主目录，返回原路径
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}
