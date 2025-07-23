#!/bin/sh

# Helper script with modular validation functions for git-send-email hook.
# Each check can be customized by setting the corresponding environment
# variable. When a command is undefined or unavailable, the check is skipped
# gracefully. This design avoids strong coupling to specific external tools
# and makes the checks easily testable in isolation.

run_cmd() {
    cmd="$1"
    shift
    if [ -n "$cmd" ] && command -v "${cmd%% *}" >/dev/null 2>&1; then
        eval "$cmd" "$@"
    else
        # Skip silently when command is not provided or not found
        return 0
    fi
}

spell_check() {
    run_cmd "${SPELL_CHECK_CMD:-}" "$@"
}

style_check() {
    run_cmd "${STYLE_CHECK_CMD:-}" "$@"
}

lint_check() {
    run_cmd "${LINT_CHECK_CMD:-}" "$@"
}

build_check() {
    run_cmd "${BUILD_CHECK_CMD:-}" "$@"
}

