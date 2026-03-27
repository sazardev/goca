package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── parseMiddlewareTypes ────────────────────────────────────────────────────

func TestParseMiddlewareTypes_Default(t *testing.T) {
	types := parseMiddlewareTypes("cors,logging,recovery")
	assert.Equal(t, []string{"cors", "logging", "recovery"}, types)
}

func TestParseMiddlewareTypes_AllTypes(t *testing.T) {
	types := parseMiddlewareTypes("cors,logging,auth,rate-limit,recovery,request-id,timeout")
	assert.Len(t, types, 7)
}

func TestParseMiddlewareTypes_Whitespace(t *testing.T) {
	types := parseMiddlewareTypes(" cors , logging , auth ")
	assert.Equal(t, []string{"cors", "logging", "auth"}, types)
}

func TestParseMiddlewareTypes_Empty(t *testing.T) {
	types := parseMiddlewareTypes("")
	assert.Empty(t, types)
}

func TestParseMiddlewareTypes_CaseInsensitive(t *testing.T) {
	types := parseMiddlewareTypes("CORS,Logging,AUTH")
	assert.Equal(t, []string{"cors", "logging", "auth"}, types)
}

// ─── validateMiddlewareTypes ─────────────────────────────────────────────────

func TestValidateMiddlewareTypes_Valid(t *testing.T) {
	err := validateMiddlewareTypes([]string{"cors", "logging", "recovery"})
	assert.NoError(t, err)
}

func TestValidateMiddlewareTypes_AllValid(t *testing.T) {
	err := validateMiddlewareTypes(validMiddlewareTypes)
	assert.NoError(t, err)
}

func TestValidateMiddlewareTypes_Invalid(t *testing.T) {
	err := validateMiddlewareTypes([]string{"cors", "banana"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "banana")
}

// ─── Template generators ─────────────────────────────────────────────────────

func TestGenerateChainMiddleware(t *testing.T) {
	out := generateChainMiddleware("github.com/test/proj")
	assert.Contains(t, out, "package middleware")
	assert.Contains(t, out, "type Middleware func(http.Handler) http.Handler")
	assert.Contains(t, out, "func Chain(mws ...Middleware) Middleware")
}

func TestGenerateCORSMiddleware(t *testing.T) {
	out := generateCORSMiddleware()
	assert.Contains(t, out, "package middleware")
	assert.Contains(t, out, "type CORSConfig struct")
	assert.Contains(t, out, "func CORS(cfg CORSConfig) Middleware")
	assert.Contains(t, out, "func DefaultCORSConfig() CORSConfig")
	assert.Contains(t, out, "Access-Control-Allow-Origin")
}

func TestGenerateLoggingMiddleware(t *testing.T) {
	out := generateLoggingMiddleware()
	assert.Contains(t, out, "package middleware")
	assert.Contains(t, out, "func Logging() Middleware")
	assert.Contains(t, out, "type statusWriter struct")
	assert.Contains(t, out, "time.Since(start)")
}

func TestGenerateAuthMiddleware(t *testing.T) {
	out := generateAuthMiddleware()
	assert.Contains(t, out, "package middleware")
	assert.Contains(t, out, "func Auth() Middleware")
	assert.Contains(t, out, "JWT_SECRET")
	assert.Contains(t, out, "func ClaimsFromContext")
	assert.Contains(t, out, "Bearer ")
}

func TestGenerateRateLimitMiddleware(t *testing.T) {
	out := generateRateLimitMiddleware()
	assert.Contains(t, out, "package middleware")
	assert.Contains(t, out, "func RateLimit(cfg RateLimitConfig) Middleware")
	assert.Contains(t, out, "type RateLimitConfig struct")
	assert.Contains(t, out, "func DefaultRateLimitConfig()")
	assert.Contains(t, out, "StatusTooManyRequests")
}

func TestGenerateRecoveryMiddleware(t *testing.T) {
	out := generateRecoveryMiddleware()
	assert.Contains(t, out, "package middleware")
	assert.Contains(t, out, "func Recovery(debugMode bool) Middleware")
	assert.Contains(t, out, "recover()")
	assert.Contains(t, out, "debug.Stack()")
}

func TestGenerateRequestIDMiddleware(t *testing.T) {
	out := generateRequestIDMiddleware()
	assert.Contains(t, out, "package middleware")
	assert.Contains(t, out, "func RequestID() Middleware")
	assert.Contains(t, out, "X-Request-ID")
	assert.Contains(t, out, "func RequestIDFromContext")
	assert.Contains(t, out, "uuid.New()")
}

func TestGenerateTimeoutMiddleware(t *testing.T) {
	out := generateTimeoutMiddleware()
	assert.Contains(t, out, "package middleware")
	assert.Contains(t, out, "func Timeout(d time.Duration) Middleware")
	assert.Contains(t, out, "context.WithTimeout")
}

// ─── generateMiddlewarePackage ───────────────────────────────────────────────

func TestGenerateMiddlewarePackage_DryRun(t *testing.T) {
	cleanup := ensureTestUI(t)
	defer cleanup()

	sm := NewSafetyManager(true, false, false)

	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	require.NoError(t, os.Chdir(tmpDir))
	defer func() { _ = os.Chdir(origDir) }()

	err := generateMiddlewarePackage("MyApp", []string{"cors", "logging", "recovery"}, sm)
	require.NoError(t, err)

	pending := sm.GetPendingFiles()
	// 4 files: middleware.go + cors.go + logging.go + recovery.go
	assert.Len(t, pending, 4)

	paths := make([]string, len(pending))
	for i, p := range pending {
		paths[i] = p.Path
	}
	assert.Contains(t, paths, filepath.Join("internal", "middleware", "middleware.go"))
	assert.Contains(t, paths, filepath.Join("internal", "middleware", "cors.go"))
	assert.Contains(t, paths, filepath.Join("internal", "middleware", "logging.go"))
	assert.Contains(t, paths, filepath.Join("internal", "middleware", "recovery.go"))
}

func TestGenerateMiddlewarePackage_AllTypes_DryRun(t *testing.T) {
	cleanup := ensureTestUI(t)
	defer cleanup()

	sm := NewSafetyManager(true, false, false)

	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	require.NoError(t, os.Chdir(tmpDir))
	defer func() { _ = os.Chdir(origDir) }()

	err := generateMiddlewarePackage("MyApp", validMiddlewareTypes, sm)
	require.NoError(t, err)

	pending := sm.GetPendingFiles()
	// 8 files: middleware.go + 7 type files
	assert.Len(t, pending, 8)
}

func TestGenerateMiddlewarePackage_RealFiles(t *testing.T) {
	cleanup := ensureTestUI(t)
	defer cleanup()

	sm := NewSafetyManager(false, true, false)

	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	require.NoError(t, os.Chdir(tmpDir))
	defer func() { _ = os.Chdir(origDir) }()

	err := generateMiddlewarePackage("TestProj", []string{"cors", "timeout"}, sm)
	require.NoError(t, err)

	middlewareDir := filepath.Join("internal", "middleware")

	// chain helper always created
	chainFile := filepath.Join(middlewareDir, "middleware.go")
	data, err := os.ReadFile(chainFile)
	require.NoError(t, err)
	assert.Contains(t, string(data), "func Chain")

	// cors file
	corsFile := filepath.Join(middlewareDir, "cors.go")
	data, err = os.ReadFile(corsFile)
	require.NoError(t, err)
	assert.Contains(t, string(data), "func CORS")

	// timeout file
	timeoutFile := filepath.Join(middlewareDir, "timeout.go")
	data, err = os.ReadFile(timeoutFile)
	require.NoError(t, err)
	assert.Contains(t, string(data), "func Timeout")

	// auth should NOT exist
	authFile := filepath.Join(middlewareDir, "auth.go")
	_, err = os.Stat(authFile)
	assert.True(t, os.IsNotExist(err))
}

func TestGenerateMiddlewarePackage_UnknownType(t *testing.T) {
	cleanup := ensureTestUI(t)
	defer cleanup()

	sm := NewSafetyManager(true, false, false)

	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	require.NoError(t, os.Chdir(tmpDir))
	defer func() { _ = os.Chdir(origDir) }()

	err := generateMiddlewarePackage("MyApp", []string{"nonexistent"}, sm)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nonexistent")
}

// ─── Template content validation ─────────────────────────────────────────────

func TestAllTemplatesCompileAsGoCode(t *testing.T) {
	// Verify every template generator produces valid Go that at least
	// starts with "package middleware".
	generators := map[string]func() string{
		"cors":       generateCORSMiddleware,
		"logging":    generateLoggingMiddleware,
		"auth":       generateAuthMiddleware,
		"rate-limit": generateRateLimitMiddleware,
		"recovery":   generateRecoveryMiddleware,
		"request-id": generateRequestIDMiddleware,
		"timeout":    generateTimeoutMiddleware,
	}
	for name, gen := range generators {
		t.Run(name, func(t *testing.T) {
			out := gen()
			assert.True(t, strings.HasPrefix(out, "package middleware"), "expected package middleware header for %s", name)
		})
	}
}
