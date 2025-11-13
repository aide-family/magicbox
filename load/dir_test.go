package load_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aide-family/magicbox/load"
)

// TestExpandHomeDir_WithTilde 测试 ExpandHomeDir 处理 ~/ 开头的路径
func TestExpandHomeDir_WithTilde(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	// 测试基本路径
	path := "~/test"
	expected := filepath.Join(home, "test")
	result := load.ExpandHomeDir(path)
	if result != expected {
		t.Errorf("ExpandHomeDir(%q) = %q, want %q", path, result, expected)
	}

	// 测试嵌套路径
	path2 := "~/test/subdir/file.txt"
	expected2 := filepath.Join(home, "test", "subdir", "file.txt")
	result2 := load.ExpandHomeDir(path2)
	if result2 != expected2 {
		t.Errorf("ExpandHomeDir(%q) = %q, want %q", path2, result2, expected2)
	}

	// 测试只有 ~/
	path3 := "~/"
	expected3 := filepath.Join(home, "")
	result3 := load.ExpandHomeDir(path3)
	if result3 != expected3 {
		t.Errorf("ExpandHomeDir(%q) = %q, want %q", path3, result3, expected3)
	}
}

// TestExpandHomeDir_WithoutTilde 测试 ExpandHomeDir 处理不以 ~/ 开头的路径
func TestExpandHomeDir_WithoutTilde(t *testing.T) {
	// 测试绝对路径
	path := "/absolute/path"
	result := load.ExpandHomeDir(path)
	if result != path {
		t.Errorf("ExpandHomeDir(%q) = %q, want %q", path, result, path)
	}

	// 测试相对路径
	path2 := "relative/path"
	result2 := load.ExpandHomeDir(path2)
	if result2 != path2 {
		t.Errorf("ExpandHomeDir(%q) = %q, want %q", path2, result2, path2)
	}

	// 测试空字符串
	path3 := ""
	result3 := load.ExpandHomeDir(path3)
	if result3 != path3 {
		t.Errorf("ExpandHomeDir(%q) = %q, want %q", path3, result3, path3)
	}

	// 测试包含 ~ 但不是开头的路径
	path4 := "path/~/test"
	result4 := load.ExpandHomeDir(path4)
	if result4 != path4 {
		t.Errorf("ExpandHomeDir(%q) = %q, want %q", path4, result4, path4)
	}

	// 测试单独的 ~（应该展开为用户主目录）
	path5 := "~"
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}
	result5 := load.ExpandHomeDir(path5)
	if result5 != home {
		t.Errorf("ExpandHomeDir(%q) = %q, want %q", path5, result5, home)
	}

	// 测试 ~ 后面直接跟字符（不是 /），应该用 home 替换 ~
	path6 := "~test"
	home2, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}
	result6 := load.ExpandHomeDir(path6)
	expected6 := filepath.Join(home2, "test")
	if result6 != expected6 {
		t.Errorf("ExpandHomeDir(%q) = %q, want %q", path6, result6, expected6)
	}
}

// TestExpandHomeDir_EdgeCases 测试边界情况
func TestExpandHomeDir_EdgeCases(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	// 测试 ~/ 后面是空字符串（实际上就是 ~/）
	path := "~/"
	result := load.ExpandHomeDir(path)
	expected := filepath.Join(home, "")
	if result != expected {
		t.Errorf("ExpandHomeDir(%q) = %q, want %q", path, result, expected)
	}

	// 测试 ~/ 后面只有一个字符
	path2 := "~/a"
	result2 := load.ExpandHomeDir(path2)
	expected2 := filepath.Join(home, "a")
	if result2 != expected2 {
		t.Errorf("ExpandHomeDir(%q) = %q, want %q", path2, result2, expected2)
	}

	// 测试包含特殊字符的路径
	path3 := "~/test with spaces/file.txt"
	result3 := load.ExpandHomeDir(path3)
	expected3 := filepath.Join(home, "test with spaces", "file.txt")
	if result3 != expected3 {
		t.Errorf("ExpandHomeDir(%q) = %q, want %q", path3, result3, expected3)
	}

	// 测试包含 .. 的路径
	path4 := "~/test/../other"
	result4 := load.ExpandHomeDir(path4)
	expected4 := filepath.Join(home, "test", "..", "other")
	if result4 != expected4 {
		t.Errorf("ExpandHomeDir(%q) = %q, want %q", path4, result4, expected4)
	}
}

// TestExpandHomeDir_PathSeparator 测试路径分隔符处理
func TestExpandHomeDir_PathSeparator(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	// 测试使用正确的路径分隔符
	path := "~/test/subdir"
	result := load.ExpandHomeDir(path)
	// filepath.Join 会根据操作系统处理路径分隔符
	expected := filepath.Join(home, "test", "subdir")
	if result != expected {
		t.Errorf("ExpandHomeDir(%q) = %q, want %q", path, result, expected)
	}

	// 测试多个连续的分隔符（filepath.Join 会处理）
	path2 := "~/test///subdir"
	result2 := load.ExpandHomeDir(path2)
	expected2 := filepath.Join(home, "test", "subdir")
	if result2 != expected2 {
		t.Errorf("ExpandHomeDir(%q) = %q, want %q", path2, result2, expected2)
	}
}

// TestExpandHomeDir_RealPath 测试实际路径展开
func TestExpandHomeDir_RealPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	// 测试常见的路径
	testCases := []struct {
		input    string
		expected string
	}{
		{"~/Documents", filepath.Join(home, "Documents")},
		{"~/Downloads/file.txt", filepath.Join(home, "Downloads", "file.txt")},
		{"~/.config", filepath.Join(home, ".config")},
		{"~/Library/Application Support", filepath.Join(home, "Library", "Application Support")},
	}

	for _, tc := range testCases {
		result := load.ExpandHomeDir(tc.input)
		if result != tc.expected {
			t.Errorf("ExpandHomeDir(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}

// BenchmarkExpandHomeDir_WithTilde 基准测试 ExpandHomeDir（带 ~/）
func BenchmarkExpandHomeDir_WithTilde(b *testing.B) {
	path := "~/test/path/to/file.txt"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = load.ExpandHomeDir(path)
	}
}

// BenchmarkExpandHomeDir_WithoutTilde 基准测试 ExpandHomeDir（不带 ~/）
func BenchmarkExpandHomeDir_WithoutTilde(b *testing.B) {
	path := "/absolute/path/to/file.txt"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = load.ExpandHomeDir(path)
	}
}

