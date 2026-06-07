#!/usr/bin/env bash
set -euo pipefail

# ──────────────────────────────────────────────
# Goca — Pre-Commit/Pre-Push Hooks Installer
# Installs: lefthook, gofumpt, gosec, nilaway,
#           staticcheck, golangci-lint
# ──────────────────────────────────────────────

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(dirname "$SCRIPT_DIR")"
cd "$REPO_DIR"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

pass() { echo -e "  ${GREEN}✓${NC} $1"; }
warn() { echo -e "  ${YELLOW}⚠${NC} $1"; }
fail() { echo -e "  ${RED}✗${NC} $1"; }
info() { echo -e "  ${CYAN}→${NC} $1"; }

echo ""
echo -e "${CYAN}═══════════════════════════════════════════${NC}"
echo -e "${CYAN}  Goca Hook Installer                      ${NC}"
echo -e "${CYAN}═══════════════════════════════════════════${NC}"
echo ""

# ── Prerequisites ──────────────────────────────

PREREQ_FAIL=0

if ! command -v go &>/dev/null; then
  fail "Go is not installed. Install Go 1.25+ first."
  PREREQ_FAIL=1
else
  GO_VERSION=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+')
  pass "Go $GO_VERSION detected"
fi

if ! command -v node &>/dev/null; then
  warn "Node.js not found — docs build will fail. Install Node 18+."
else
  NODE_VERSION=$(node --version | cut -d'v' -f2 | cut -d'.' -f1)
  pass "Node.js v$(node --version) detected"
fi

if ! command -v npm &>/dev/null; then
  warn "npm not found — docs build will fail."
else
  pass "npm $(npm --version) detected"
fi

if [ "$PREREQ_FAIL" -eq 1 ]; then
  echo ""
  fail "Prerequisites missing. Aborting."
  exit 1
fi

echo ""

# ── Install Go Tools ───────────────────────────

install_go_tool() {
  local cmd=$1
  local pkg=$2
  local label=${3:-$cmd}

  if command -v "$cmd" &>/dev/null; then
    pass "$label already installed ($($cmd version 2>&1 | head -1))"
  else
    info "Installing $label..."
    go install "$pkg@latest"
    if command -v "$cmd" &>/dev/null; then
      pass "$label installed"
    else
      warn "$label installed but not in PATH — add \$(go env GOPATH)/bin to your PATH"
    fi
  fi
}

install_go_tool "lefthook" "github.com/evilmartians/lefthook@latest" "lefthook"
install_go_tool "gofumpt" "mvdan.cc/gofumpt@latest" "gofumpt"
install_go_tool "gosec" "github.com/securego/gosec/v2/cmd/gosec@latest" "gosec"
install_go_tool "nilaway" "go.uber.org/nilaway/cmd/nilaway@latest" "nilaway"
install_go_tool "staticcheck" "honnef.co/go/tools/cmd/staticcheck@latest" "staticcheck"

# ── golangci-lint (separate — different install pattern) ──

if command -v golangci-lint &>/dev/null; then
  pass "golangci-lint already installed ($(golangci-lint version 2>&1 | head -1))"
else
  info "Installing golangci-lint..."
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOPATH)/bin"
  if command -v golangci-lint &>/dev/null; then
    pass "golangci-lint installed"
  else
    warn "golangci-lint installation may need manual setup"
  fi
fi

echo ""
info "Verifying tools..."

TOOLS=(lefthook gofumpt gosec nilaway staticcheck golangci-lint go node npm)
ALL_OK=0
for tool in "${TOOLS[@]}"; do
  if command -v "$tool" &>/dev/null; then
    pass "$tool found"
  else
    warn "$tool not found (may be optional)"
  fi
done

echo ""

# ── Install Docs Dependencies ──────────────────

if [ -f "docs/package.json" ]; then
  info "Installing docs dependencies..."
  (cd docs && npm ci 2>&1 | tail -1)
  pass "Docs dependencies installed"
else
  warn "docs/package.json not found — skipping docs deps"
fi

echo ""

# ── Install Lefthook Hooks ─────────────────────

info "Installing git hooks via lefthook..."
lefthook install
pass "Git hooks installed"

echo ""
echo -e "${CYAN}───────────────────────────────────────────────${NC}"
echo -e "${GREEN}  ✅ All hooks installed successfully!${NC}"
echo -e "${CYAN}───────────────────────────────────────────────${NC}"
echo ""
echo -e "  ${YELLOW}pre-commit${NC} runs: gofumpt + goimports + golangci-lint"
echo -e "         + staticcheck + gosec + nilaway + tests + build + docs"
echo ""
echo -e "  ${YELLOW}pre-push${NC}   runs: vet + golangci-lint + staticcheck"
echo -e "         + gosec + nilaway + full tests + full build + docs"
echo ""
echo -e "  ${CYAN}Skip hooks temporarily:${NC} git commit --no-verify"
echo -e "  ${CYAN}Reinstall hooks:${NC}       lefthook install -f"
echo -e "  ${CYAN}Remove hooks:${NC}          lefthook uninstall"
echo ""
