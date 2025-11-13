package safety_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/magicbox/log"
	"github.com/aide-family/magicbox/safety"
)

// mockLogger 是一个 mock logger 实现，用于测试
type mockLogger struct {
	buf *bytes.Buffer
}

// Log 实现 log.Logger 接口
func (m *mockLogger) Log(level klog.Level, keyvals ...interface{}) error {
	if m.buf == nil {
		return nil
	}
	// 简单的日志格式：level + keyvals
	var parts []string
	parts = append(parts, level.String())
	for _, kv := range keyvals {
		parts = append(parts, stringify(kv))
	}
	m.buf.WriteString(strings.Join(parts, " ") + "\n")
	return nil
}

// stringify 将值转换为字符串
func stringify(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case error:
		return val.Error()
	default:
		return ""
	}
}

// newMockLogger 创建一个新的 mock logger，将日志写入 buffer
func newMockLogger(buf *bytes.Buffer) log.Interface {
	return &mockLogger{buf: buf}
}

// newNopMockLogger 创建一个不输出任何内容的 mock logger
func newNopMockLogger() log.Interface {
	return &mockLogger{buf: nil}
}

// TestGo_Success 测试 Go 函数正常执行的情况
func TestGo_Success(t *testing.T) {
	// 创建一个 buffer 来捕获日志输出
	var buf bytes.Buffer
	logger := newMockLogger(&buf)

	ctx := context.Background()
	executed := make(chan bool, 1)
	var executedOnce sync.Once

	// 创建一个成功执行的函数
	f := func(context.Context) error {
		executedOnce.Do(func() {
			executed <- true
		})
		return nil
	}

	// 调用 Go 函数
	safety.Go(ctx, "test-success", f, logger)

	// 等待函数执行完成
	select {
	case <-executed:
		// 函数已执行
	case <-time.After(2 * time.Second):
		t.Fatal("Function did not execute within timeout")
	}

	// 等待日志记录完成（defer 中的日志）
	time.Sleep(100 * time.Millisecond)

	// 验证日志输出包含 completed 信息
	logOutput := buf.String()
	if !strings.Contains(logOutput, "completed") {
		t.Errorf("Expected log to contain 'completed', got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "cost") {
		t.Errorf("Expected log to contain 'cost', got: %s", logOutput)
	}
}

// TestGo_Error 测试 Go 函数处理错误的情况
func TestGo_Error(t *testing.T) {
	// 创建一个 buffer 来捕获日志输出
	var buf bytes.Buffer
	logger := newMockLogger(&buf)

	ctx := context.Background()
	testErr := errors.New("test error")
	executed := make(chan bool, 1)
	var executedOnce sync.Once

	// 创建一个返回错误的函数
	f := func(context.Context) error {
		executedOnce.Do(func() {
			executed <- true
		})
		return testErr
	}

	// 调用 Go 函数
	safety.Go(ctx, "test-error", f, logger)

	// 等待函数执行完成
	select {
	case <-executed:
		// 函数已执行
	case <-time.After(2 * time.Second):
		t.Fatal("Function did not execute within timeout")
	}

	// 等待日志记录完成
	time.Sleep(100 * time.Millisecond)

	// 验证日志输出包含错误信息
	logOutput := buf.String()
	if !strings.Contains(logOutput, "run error") {
		t.Errorf("Expected log to contain 'run error', got: %s", logOutput)
	}
	if !strings.Contains(logOutput, testErr.Error()) {
		t.Errorf("Expected log to contain error message '%s', got: %s", testErr.Error(), logOutput)
	}
	if !strings.Contains(logOutput, "completed") {
		t.Errorf("Expected log to contain 'completed', got: %s", logOutput)
	}
}

// TestGo_Panic 测试 Go 函数处理 panic 的情况
func TestGo_Panic(t *testing.T) {
	// 创建一个 buffer 来捕获日志输出
	var buf bytes.Buffer
	logger := newMockLogger(&buf)

	ctx := context.Background()
	panicMsg := "test panic"
	executed := make(chan bool, 1)
	var executedOnce sync.Once

	// 创建一个会 panic 的函数
	f := func(context.Context) error {
		executedOnce.Do(func() {
			executed <- true
		})
		panic(panicMsg)
	}

	// 调用 Go 函数
	safety.Go(ctx, "test-panic", f, logger)

	// 等待函数执行完成
	select {
	case <-executed:
		// 函数已执行（虽然会 panic）
	case <-time.After(2 * time.Second):
		t.Fatal("Function did not execute within timeout")
	}

	// 等待 panic 被捕获并记录日志
	time.Sleep(100 * time.Millisecond)

	// 验证日志输出包含 panic 信息
	logOutput := buf.String()
	if !strings.Contains(logOutput, "panic") {
		t.Errorf("Expected log to contain 'panic', got: %s", logOutput)
	}
	if !strings.Contains(logOutput, panicMsg) {
		t.Errorf("Expected log to contain panic message '%s', got: %s", panicMsg, logOutput)
	}
	if !strings.Contains(logOutput, "completed") {
		t.Errorf("Expected log to contain 'completed', got: %s", logOutput)
	}
}

// TestGo_ContextPassed 测试 Go 函数正确传递 context
func TestGo_ContextPassed(t *testing.T) {
	// 创建一个 buffer 来捕获日志输出
	var buf bytes.Buffer
	logger := newMockLogger(&buf)

	// 创建一个带值的 context
	type key string
	testKey := key("test-key")
	testValue := "test-value"
	ctx := context.WithValue(context.Background(), testKey, testValue)

	executed := make(chan bool, 1)
	var executedOnce sync.Once
	var receivedValue interface{}

	// 创建一个验证 context 的函数
	f := func(c context.Context) error {
		executedOnce.Do(func() {
			receivedValue = c.Value(testKey)
			executed <- true
		})
		return nil
	}

	// 调用 Go 函数
	safety.Go(ctx, "test-context", f, logger)

	// 等待函数执行完成
	select {
	case <-executed:
		// 函数已执行
	case <-time.After(2 * time.Second):
		t.Fatal("Function did not execute within timeout")
	}

	// 验证 context 值被正确传递
	if receivedValue != testValue {
		t.Errorf("Expected context value %v, got %v", testValue, receivedValue)
	}
}

// TestGo_ConcurrentExecution 测试 Go 函数并发执行
func TestGo_ConcurrentExecution(t *testing.T) {
	// 创建一个 buffer 来捕获日志输出（虽然并发时顺序不确定）
	var buf bytes.Buffer
	logger := newMockLogger(&buf)

	ctx := context.Background()
	var wg sync.WaitGroup
	const numGoroutines = 10

	// 创建多个函数并发执行
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		f := func(context.Context) error {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			return nil
		}
		safety.Go(ctx, "test-concurrent", f, logger)
	}

	// 等待所有 goroutine 完成
	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		// 所有 goroutine 已完成
	case <-time.After(5 * time.Second):
		t.Fatal("Not all goroutines completed within timeout")
	}

	// 等待日志记录完成
	time.Sleep(200 * time.Millisecond)

	// 验证至少有一些日志输出
	logOutput := buf.String()
	if len(logOutput) == 0 {
		t.Error("Expected some log output, got empty")
	}
}

// TestGo_WithCancelledContext 测试 Go 函数处理已取消的 context
func TestGo_WithCancelledContext(t *testing.T) {
	// 创建一个 buffer 来捕获日志输出
	var buf bytes.Buffer
	logger := newMockLogger(&buf)

	// 创建一个已取消的 context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	executed := make(chan bool, 1)
	var executedOnce sync.Once

	// 创建一个检查 context 状态的函数
	f := func(c context.Context) error {
		executedOnce.Do(func() {
			executed <- true
		})
		// 即使 context 已取消，函数仍应执行
		if c.Err() != nil {
			// Context 已取消，这是预期的
		}
		return nil
	}

	// 调用 Go 函数
	safety.Go(ctx, "test-cancelled-context", f, logger)

	// 等待函数执行完成
	select {
	case <-executed:
		// 函数已执行
	case <-time.After(2 * time.Second):
		t.Fatal("Function did not execute within timeout")
	}

	// 等待日志记录完成
	time.Sleep(100 * time.Millisecond)

	// 验证日志输出包含 completed 信息
	logOutput := buf.String()
	if !strings.Contains(logOutput, "completed") {
		t.Errorf("Expected log to contain 'completed', got: %s", logOutput)
	}
}

// TestGo_EmptyName 测试 Go 函数使用空名称
func TestGo_EmptyName(t *testing.T) {
	// 创建一个 buffer 来捕获日志输出
	var buf bytes.Buffer
	logger := newMockLogger(&buf)

	ctx := context.Background()
	executed := make(chan bool, 1)
	var executedOnce sync.Once

	// 创建一个成功执行的函数
	f := func(context.Context) error {
		executedOnce.Do(func() {
			executed <- true
		})
		return nil
	}

	// 调用 Go 函数，使用空名称
	safety.Go(ctx, "", f, logger)

	// 等待函数执行完成
	select {
	case <-executed:
		// 函数已执行
	case <-time.After(2 * time.Second):
		t.Fatal("Function did not execute within timeout")
	}

	// 等待日志记录完成
	time.Sleep(100 * time.Millisecond)

	// 验证日志输出包含 completed 信息
	logOutput := buf.String()
	if !strings.Contains(logOutput, "completed") {
		t.Errorf("Expected log to contain 'completed', got: %s", logOutput)
	}
}

// TestGo_WithNopLogger 测试 Go 函数使用 NopLogger
func TestGo_WithNopLogger(t *testing.T) {
	// 使用 nop mock logger（不输出任何内容）
	logger := newNopMockLogger()

	ctx := context.Background()
	executed := make(chan bool, 1)
	var executedOnce sync.Once

	// 创建一个成功执行的函数
	f := func(context.Context) error {
		executedOnce.Do(func() {
			executed <- true
		})
		return nil
	}

	// 调用 Go 函数
	safety.Go(ctx, "test-nop-logger", f, logger)

	// 等待函数执行完成
	select {
	case <-executed:
		// 函数已执行
	case <-time.After(2 * time.Second):
		t.Fatal("Function did not execute within timeout")
	}

	// 等待日志记录完成（虽然不会输出）
	time.Sleep(100 * time.Millisecond)

	// 如果执行到这里，说明没有 panic，测试通过
}

// BenchmarkGo 基准测试 Go 函数
func BenchmarkGo(b *testing.B) {
	logger := newNopMockLogger()
	ctx := context.Background()
	f := func(context.Context) error {
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(1)
		func() {
			defer wg.Done()
			safety.Go(ctx, "benchmark", f, logger)
		}()
		wg.Wait()
		// 等待 goroutine 完成
		time.Sleep(10 * time.Millisecond)
	}
}
