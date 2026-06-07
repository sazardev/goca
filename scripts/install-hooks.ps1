# Goca — Pre-Commit/Pre-Push Hooks Installer (Windows)
# Installs: lefthook, gofumpt, gosec, nilaway,
#           staticcheck, golangci-lint

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

function Write-Pass  { Write-Host "  ✓ $args" -ForegroundColor Green }
function Write-Warn  { Write-Host "  ⚠ $args" -ForegroundColor Yellow }
function Write-Fail  { Write-Host "  ✗ $args" -ForegroundColor Red }
function Write-Info  { Write-Host "  → $args" -ForegroundColor Cyan }

$repoDir = Split-Path -Parent $PSScriptRoot
Set-Location $repoDir

Write-Host ""
Write-Host "═══════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "  Goca Hook Installer (Windows)" -ForegroundColor Cyan
Write-Host "═══════════════════════════════════════════" -ForegroundColor Cyan
Write-Host ""

# ── Prerequisites ──

$prereqFail = $false

$goCmd = Get-Command "go" -ErrorAction SilentlyContinue
if (-not $goCmd) {
  Write-Fail "Go is not installed. Install Go 1.25+ first."
  $prereqFail = $true
} else {
  $goVersion = go version
  Write-Pass "Go detected: $goVersion"
}

$nodeCmd = Get-Command "node" -ErrorAction SilentlyContinue
if (-not $nodeCmd) {
  Write-Warn "Node.js not found — docs build will fail. Install Node 18+."
} else {
  Write-Pass "Node.js $(node --version) detected"
}

$npmCmd = Get-Command "npm" -ErrorAction SilentlyContinue
if (-not $npmCmd) {
  Write-Warn "npm not found — docs build will fail."
} else {
  Write-Pass "npm $(npm --version) detected"
}

if ($prereqFail) {
  Write-Host ""
  Write-Fail "Prerequisites missing. Aborting."
  exit 1
}

Write-Host ""

# ── Install Go Tools ──

function Install-GoTool($cmd, $pkg, $label) {
  $existing = Get-Command $cmd -ErrorAction SilentlyContinue
  if ($existing) {
    Write-Pass "$label already installed"
  } else {
    Write-Info "Installing $label..."
    go install "$pkg@latest"
    $gopath = go env GOPATH
    $binPath = Join-Path $gopath "bin"
    $toolPath = Join-Path $binPath $cmd
    if (Test-Path $toolPath) {
      Write-Pass "$label installed"
    } else {
      Write-Warn "$label installed but not in PATH — add $binPath to your PATH"
    }
  }
}

Install-GoTool "lefthook"    "github.com/evilmartians/lefthook@latest"              "lefthook"
Install-GoTool "gofumpt"     "mvdan.cc/gofumpt@latest"                               "gofumpt"
Install-GoTool "gosec"       "github.com/securego/gosec/v2/cmd/gosec@latest"         "gosec"
Install-GoTool "nilaway"     "go.uber.org/nilaway/cmd/nilaway@latest"                "nilaway"
Install-GoTool "staticcheck" "honnef.co/go/tools/cmd/staticcheck@latest"            "staticcheck"

# ── golangci-lint ──

$gci = Get-Command "golangci-lint" -ErrorAction SilentlyContinue
if ($gci) {
  Write-Pass "golangci-lint already installed"
} else {
  Write-Info "Installing golangci-lint..."
  $gopath = go env GOPATH
  $binPath = Join-Path $gopath "bin"
  $installUrl = "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh"
  $installer = Join-Path $env:TEMP "install-golangci.sh"
  Invoke-WebRequest -Uri $installUrl -OutFile $installer -UseBasicParsing
  bash $installer -b $binPath
  if (Test-Path (Join-Path $binPath "golangci-lint")) {
    Write-Pass "golangci-lint installed"
  } else {
    Write-Warn "golangci-lint installation may need manual setup"
  }
}

Write-Host ""
Write-Info "Verifying tools..."

$tools = @("lefthook", "gofumpt", "gosec", "nilaway", "staticcheck", "golangci-lint", "go", "node", "npm")
foreach ($tool in $tools) {
  $found = Get-Command $tool -ErrorAction SilentlyContinue
  if ($found) {
    Write-Pass "$tool found"
  } else {
    Write-Warn "$tool not found (may be optional)"
  }
}

Write-Host ""

# ── Install Docs Dependencies ──

if (Test-Path "docs/package.json") {
  Write-Info "Installing docs dependencies..."
  Push-Location docs
  npm ci 2>&1 | Select-Object -Last 1
  Pop-Location
  Write-Pass "Docs dependencies installed"
} else {
  Write-Warn "docs/package.json not found — skipping docs deps"
}

Write-Host ""

# ── Install Lefthook Hooks ──

Write-Info "Installing git hooks via lefthook..."
lefthook install
Write-Pass "Git hooks installed"

Write-Host ""
Write-Host "───────────────────────────────────────────────" -ForegroundColor Cyan
Write-Host "  ✅ All hooks installed successfully!" -ForegroundColor Green
Write-Host "───────────────────────────────────────────────" -ForegroundColor Cyan
Write-Host ""
Write-Host "  pre-commit runs: gofumpt + goimports + golangci-lint" -ForegroundColor Yellow
Write-Host "         + staticcheck + gosec + nilaway + tests + build + docs"
Write-Host ""
Write-Host "  pre-push   runs: vet + golangci-lint + staticcheck" -ForegroundColor Yellow
Write-Host "         + gosec + nilaway + full tests + full build + docs"
Write-Host ""
Write-Host "  Skip hooks temporarily: git commit --no-verify" -ForegroundColor Cyan
Write-Host "  Reinstall hooks:       lefthook install -f" -ForegroundColor Cyan
Write-Host "  Remove hooks:          lefthook uninstall" -ForegroundColor Cyan
Write-Host ""
