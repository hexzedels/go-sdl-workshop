#!/usr/bin/env bash
# lib.sh — shared helpers for per-challenge verify scripts.
# Each verify-fix/<id>.sh sources this and calls run_fix_check.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"

# shellcheck source=../lib/flags.sh
source "$SCRIPT_DIR/../lib/flags.sh"

# run_fix_check <short-id> <challenge-name> <go-test-regex>
# Runs the static fix check, then reveals the fix flag if it passes.
run_fix_check() {
    local short_id="$1"
    local challenge_name="$2"
    local test_regex="$3"

    cd "$ROOT_DIR"

    echo "=== Verifying fix: $short_id ==="
    echo ""
    if ! go test -v -run "$test_regex" ./verify/; then
        echo ""
        echo "Fix check for $short_id failed. See errors above."
        return 1
    fi

    local flag
    if ! flag=$(fix_flag_for "$short_id"); then
        echo ""
        echo "lib.sh: no flag registered for challenge id '$short_id' in scripts/lib/flags.sh" >&2
        return 1
    fi

    echo ""
    echo "=========================================="
    echo "Fix verified. Submit this flag to CTFd:"
    echo "  $short_id ($challenge_name) → $flag"
    echo "=========================================="
}
