#!/usr/bin/env bash
# verify-setup.sh — Pre-workshop environment check
# Run this before the workshop to ensure all tools are installed.

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASS=0
FAIL=0
WARN=0

version_of() {
    case "$1" in
        go) go version 2>&1 | head -1 ;;
        *) "$1" --version 2>&1 | head -1 ;;
    esac
}

check() {
    local name="$1"
    local cmd="$2"

    if command -v "$cmd" &> /dev/null; then
        local version
        version=$(version_of "$cmd" || echo "unknown")
        echo -e "  ${GREEN}✓${NC} $name: $version"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}✗${NC} $name: not found"
        FAIL=$((FAIL + 1))
    fi
}

warn_check() {
    local name="$1"
    local cmd="$2"

    if command -v "$cmd" &> /dev/null; then
        echo -e "  ${GREEN}✓${NC} $name: found"
        PASS=$((PASS + 1))
    else
        echo -e "  ${YELLOW}?${NC} $name: not found (optional)"
        WARN=$((WARN + 1))
    fi
}

echo "GoSDLWorkshop — Environment Check"
echo "=================================="
echo ""
echo "Required tools:"
check "Go" "go"
check "Git" "git"
check "curl" "curl"
check "jq" "jq"

echo ""
echo "Container tools:"
check "Docker" "docker"
warn_check "Docker Compose" "docker"

echo ""
echo "Recommended tools (installed during workshop):"
warn_check "gosec" "gosec"
warn_check "govulncheck" "govulncheck"
warn_check "trivy" "trivy"
warn_check "semgrep" "semgrep"

echo ""
echo "Checking Go version..."
GO_VERSION=$(go version 2>/dev/null | awk '{for (i=1;i<=NF;i++) if ($i ~ /^go[0-9]/) {sub(/^go/, "", $i); print $i; exit}}')
GO_MAJOR=$(echo "${GO_VERSION:-0.0}" | cut -d. -f1)
GO_MINOR=$(echo "${GO_VERSION:-0.0}" | cut -d. -f2)
if [ "${GO_MAJOR:-0}" -ge 1 ] 2>/dev/null && [ "${GO_MINOR:-0}" -ge 22 ] 2>/dev/null; then
    echo -e "  ${GREEN}✓${NC} Go $GO_VERSION (>= 1.22 required)"
else
    echo -e "  ${RED}✗${NC} Go ${GO_VERSION:-unknown} (>= 1.22 required)"
    FAIL=$((FAIL + 1))
fi

echo ""
echo "Checking project builds..."
if go build ./... 2>/dev/null; then
    echo -e "  ${GREEN}✓${NC} go build ./... succeeded"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}✗${NC} go build ./... failed"
    FAIL=$((FAIL + 1))
fi

echo ""
echo "=================================="
echo -e "Results: ${GREEN}$PASS passed${NC}, ${RED}$FAIL failed${NC}, ${YELLOW}$WARN optional missing${NC}"

if [ $FAIL -gt 0 ]; then
    echo -e "${RED}Please install missing required tools before the workshop.${NC}"
    exit 1
else
    echo -e "${GREEN}Environment is ready for the workshop!${NC}"
fi
