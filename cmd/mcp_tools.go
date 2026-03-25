package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// runGocaSubcommand invokes the current goca binary with the given args and
// returns its combined stdout/stderr output. Using os.Args[0] ensures the
// running binary is always in sync with this MCP server instance.
func runGocaSubcommand(ctx context.Context, args []string) (string, error) {
	binary, err := os.Executable()
	if err != nil {
		binary = os.Args[0]
	}

	// Safety: validate that every argument is a safe string before passing
	// it to exec.Command. Each element must not contain shell meta-characters.
	for _, a := range args {
		if strings.ContainsAny(a, "`$|;&<>(){}\\") {
			return "", fmt.Errorf("unsafe argument rejected: %q", a)
		}
	}

	// Always suppress interactive prompts when called from MCP.
	safeArgs := append([]string{"--no-interactive"}, args...)

	var out bytes.Buffer
	//nolint:gosec // binary comes from os.Executable() — trusted path
	cmd := exec.CommandContext(ctx, binary, safeArgs...)
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		// Include captured output so the LLM gets the real error message.
		return "", fmt.Errorf("%w\n%s", err, out.String())
	}

	return out.String(), nil
}

// mcpText returns a successful tool result containing text.
func mcpText(text string) *mcp.CallToolResult {
	return mcp.NewToolResultText(text)
}

// mcpErr returns an error tool result.
func mcpErr(err error) *mcp.CallToolResult {
	return mcp.NewToolResultErrorFromErr("goca error", err)
}

// appendIfTrue appends flag and value to args when condition is true.
func appendIfTrue(args []string, condition bool, flag string) []string {
	if condition {
		return append(args, flag)
	}
	return args
}

// appendIfSet appends flag and value to args when value is non-empty.
func appendIfSet(args []string, value, flag string) []string {
	if value != "" {
		return append(args, flag, value)
	}
	return args
}
