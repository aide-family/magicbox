package load

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandHomeDir expand the ~ symbol in the path to the user's home directory
func ExpandHomeDir(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			// 如果无法获取用户主目录，返回原路径
			return path
		}
		// 移除 ~ 并用 filepath.Join 拼接
		return filepath.Join(home, path[1:])
	}
	return path
}
